package oauth

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"time"
)

const (
	microsoftLoginURL = `https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&scope=%s&response_type=code&redirect_uri=%s`
)

var microsoftScope = "offline_access ServiceActivity-OneDrive.Read.All"

type (
	// MicrosoftConfig  declare configuration for Microsoft OAuth service
	MicrosoftConfig struct {
		ClientID string `env:"CLIENT_ID,required"`
	}
	// Microsoft implement OAuth authorization for microsoft`s services
	Microsoft struct {
		cfg        *MicrosoftConfig
		httpClient *resty.Client
	}
)

func (m *Microsoft) GetAuthURL(ctx context.Context, redirectUrl, state string) (string, error) {
	return fmt.Sprintf(microsoftLoginURL, m.cfg.ClientID, microsoftScope, redirectUrl), nil
}

func (m *Microsoft) ExchangeCode(ctx context.Context, code string) (*IdentityResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Microsoft) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	//TODO implement me
	panic("implement me")
}

func NewMicrosoft(cfg *MicrosoftConfig) Proxy {
	return &Microsoft{
		cfg:        cfg,
		httpClient: resty.New().SetTimeout(time.Minute),
	}
}
