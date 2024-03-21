package oauth

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"time"
)

const (
	StateNameGoogle   = "state"
	CodeNameGoogle    = "code"
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type GoogleLoginResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
type googleProvider struct {
	config     *oauth2.Config
	httpClient *resty.Client
}

// NewGoogleProvider create new implementation of google oAuth client
func NewGoogleProvider(cfg *Config, redirectURL string) Proxy {
	return &googleProvider{
		httpClient: resty.New().SetTimeout(time.Minute),
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
	}
}

func (g *googleProvider) Login(ctx context.Context) (*SignInConfig, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, utils.Internal.Wrap(err, "can not generate state")
	}
	state := base64.URLEncoding.EncodeToString(b)
	config := &SignInConfig{
		State:           state,
		StateCookieName: CodeNameGoogle,
		URL:             g.config.AuthCodeURL(state),
	}
	return config, nil
}

func (g *googleProvider) Callback(ctx context.Context, code string) (*model.JwtClaim, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, utils.Internal.Wrapf(err, "code exchange wrong: %s", err.Error())
	}
	response, err := g.httpClient.NewRequest().Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil || response.IsError() {
		return nil, utils.Internal.Newf("failed getting user info: %v [%v]", err.Error(), response.Error())
	}

	contents := response.Body()

	var data GoogleLoginResponse
	if err = json.Unmarshal(contents, &data); err != nil {
		return nil, utils.Internal.Wrapf(err, "can not marshal google response")
	}
	return &model.JwtClaim{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserName:     data.Email,
		TenantID:     uuid.New().String(),
	}, nil
}

func (g *googleProvider) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	//TODO implement me
	panic("implement me")
}
