package oauth

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"golang.org/x/oauth2"
)

type (
	SignInConfig struct {
		State           string
		StateCookieName string
		URL             string
	}

	Config struct {
		GoogleClientID string `env:"GOOGLE_CLIENT_ID"`
		GoogleSecret   string `env:"GOOGLE_SECRET"`
	}
	Proxy interface {
		Login(ctx context.Context) (*SignInConfig, error)
		Callback(ctx context.Context, code string) (*model.JwtClaim, error)
		RefreshToken(token *oauth2.Token) (*oauth2.Token, error)
	}
)
