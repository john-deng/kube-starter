package oidc

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strings"
)

const (
	Profile = "oidc"
)

type configuration struct {
	at.AutoConfiguration
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Token
type Token struct {
	at.ContextAware

	Context context.Context `json:"context"`
	Data    string          `json:"data"`
	Claims  *Claims         `json:"claims"`
}

// Token
func (c *configuration) Token(ctx context.Context) (token *Token) {
	token = new(Token)

	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		bearerToken = ctx.URLParam("token")
	}
	token.Data = strings.Replace(bearerToken, "Bearer ", "", -1)
	var err error
	token.Claims, err = DecodeWithoutVerify(token.Data)
	if err != nil {
		log.Error(err)
	}
	return
}
