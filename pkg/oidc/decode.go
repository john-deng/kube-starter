// Package jwt provides JWT manipulations.
// See https://tools.ietf.org/html/rfc7519#section-4.1.3
package oidc

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt"
	"github.com/hidevopsio/hiboot/pkg/log"
	"gopkg.in/resty.v1"
	"math/big"
	"strings"
	"time"
)

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// DecodeWithoutVerify decodes the JWT string and returns the claims.
// Note that this method does not verify the signature and always trust it.
func DecodeWithoutVerify(s string) (c *Claims, err error) {
	payload, err := DecodePayloadAsRawJSON(s)
	if err != nil {
		return nil, fmt.Errorf("could not decode the payload: %w", err)
	}
	var claims struct {
		Issuer    string `json:"iss,omitempty"`
		Subject   string `json:"sub,omitempty"`
		Name      string `json:"name,omitempty"`
		Username  string `json:"preferred_username,omitempty"`
		Email     string `json:"email,omitempty"`
		ExpiresAt int64  `json:"exp,omitempty"`
	}
	if err := json.NewDecoder(bytes.NewReader(payload)).Decode(&claims); err != nil {
		return nil, fmt.Errorf("could not decode the json of token: %w", err)
	}

	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, payload, "", "  "); err != nil {
		return nil, fmt.Errorf("could not indent the json of token: %w", err)
	}
	cls := &Claims{
		Issuer:   claims.Issuer,
		Subject:  claims.Subject,
		Name:     claims.Name,
		Username: claims.Username,
		Email:    claims.Email,
		Expiry:   time.Unix(claims.ExpiresAt, 0),
		Pretty:   prettyJson.String(),
	}

	// fill username as the value of name if it is empty
	if cls.Username == "" {
		cls.Username = cls.Name
	}

	return cls, nil
}

// DecodePayloadAsPrettyJSON decodes the JWT string and returns the pretty JSON string.
func DecodePayloadAsPrettyJSON(s string) (string, error) {
	payload, err := DecodePayloadAsRawJSON(s)
	if err != nil {
		return "", fmt.Errorf("could not decode the payload: %w", err)
	}
	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, payload, "", "  "); err != nil {
		return "", fmt.Errorf("could not indent the json of token: %w", err)
	}
	return prettyJson.String(), nil
}

// DecodePayloadAsRawJSON extracts the payload and returns the raw JSON.
func DecodePayloadAsRawJSON(s string) ([]byte, error) {
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("wants %d segments but got %d segments", 3, len(parts))
	}
	payloadJSON, err := decodePayload(parts[1])
	if err != nil {
		return nil, fmt.Errorf("could not decode the payload: %w", err)
	}
	return payloadJSON, nil
}

func decodePayload(payload string) ([]byte, error) {
	b, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("invalid base64: %w", err)
	}
	return b, nil
}

// parsePublicKey converts a JWK to an RSA public key
func parsePublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes := decodeBase64URL(jwk.N)
	eBytes := decodeBase64URL(jwk.E)

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Uint64())

	return &rsa.PublicKey{N: n, E: e}, nil
}

// decodeBase64URL decodes a base64 URL encoded string
func decodeBase64URL(s string) []byte {
	data, _ := base64.RawURLEncoding.DecodeString(s)
	return data
}

// verifyToken parses and verifies a JWT token using the JWK set
func verifyToken(tokenString string, jwkSet JWKSet) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			for _, key := range jwkSet.Keys {
				if key.Kid == kid {
					return parsePublicKey(key)
				}
			}
		}
		return nil, fmt.Errorf("unable to find appropriate key")
	})
	return token, err
}

// createKeySet creates an OIDC KeySet from a JWKSet
func createKeySet(jwkSet JWKSet) oidc.KeySet {
	keys := []crypto.PublicKey{}
	for _, jwk := range jwkSet.Keys {
		if key, err := parsePublicKey(jwk); err == nil {
			keys = append(keys, key)
		}
	}
	return &oidc.StaticKeySet{PublicKeys: keys}
}

// fetchJWKSet retrieves the JWK set from the issuer URL
func fetchJWKSet(issuerURL string) (JWKSet, error) {
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("%s/keys", issuerURL))
	if err != nil {
		return JWKSet{}, err
	}

	var jwkSet JWKSet
	if err := json.Unmarshal(resp.Body(), &jwkSet); err != nil {
		return JWKSet{}, err
	}

	return jwkSet, nil
}

// verifyOIDCToken verifies an OIDC token using the issuer's JWK set
func verifyOIDCToken(issuerURL, tokenString string) (err error) {
	// Fetch the JWK set from the issuer
	jwkSet, err := fetchJWKSet(issuerURL)
	if err != nil {
		log.Errorf("Failed to fetch JWK Set: %v", err)
		return
	}

	// Verify the token
	token, err := verifyToken(tokenString, jwkSet)
	if err != nil {
		log.Errorf("Token %v verification failed: %v", token, err)
		return
	}

	// Validate the token using OIDC
	keySet := createKeySet(jwkSet)
	verifier := oidc.NewVerifier(strings.TrimSuffix(issuerURL, "/"), keySet, &oidc.Config{
		SkipClientIDCheck: true,
	})

	_, err = verifier.Verify(context.Background(), tokenString)
	if err != nil {
		log.Errorf("OIDC Token verification failed: %v", err)
	} else {
		log.Errorf("OIDC Token verification succeeded!")
	}
	return
}
