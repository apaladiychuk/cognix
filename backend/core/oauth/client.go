package oauth

import (
	"context"
	"golang.org/x/oauth2"
)

const (
	LoginState  = "login"
	SignUpState = "signUp"
	InviteState = "invite"
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
		Login(ctx context.Context, redirectUrl, state string) (string, error)
		Callback(ctx context.Context, code string) (*IdentityResponse, error)
		RefreshToken(token *oauth2.Token) (*oauth2.Token, error)
	}
)
