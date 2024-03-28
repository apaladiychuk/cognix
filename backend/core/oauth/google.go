package oauth

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
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
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
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
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"},
		},
	}
}

func (g *googleProvider) Login(ctx context.Context, state string) (*SignInConfig, error) {
	config := &SignInConfig{
		State:           state,
		StateCookieName: CodeNameGoogle,
		URL:             g.config.AuthCodeURL(state),
	}
	return config, nil
}

func (g *googleProvider) Callback(ctx context.Context, code string) (*security.Identity, error) {
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
	loginResponse := &security.Identity{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: &model.User{
			UserName:   data.Email,
			FirstName:  data.GivenName,
			LastName:   data.FamilyName,
			ExternalID: data.ID,
			Roles:      nil,
		},
	}
	if loginResponse.User.FirstName == "" {
		loginResponse.User.FirstName = data.Name
	}
	return loginResponse, nil
}

func (g *googleProvider) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	//TODO implement me
	panic("implement me")
}
