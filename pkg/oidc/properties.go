package oidc

import (
	"github.com/hidevopsio/hiboot/pkg/at"
)

// Properties the operator properties
type Properties struct {
	at.ConfigurationProperties `value:"oidc"`

	Verify     bool     `json:"verify"`
	PublicKeys []string `json:"publicKeys"`
}
