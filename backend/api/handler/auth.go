package handler

import (
	"cognix.ch/api/v2/bll"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler  handles authentication endpoints
type AuthHandler struct {
	oauthClient oauth.Proxy
	jwtService  security.JWTService
	authBL      bll.AuthBL
}

func NewAuthHandler(oauthClient oauth.Proxy,
	jwtService security.JWTService,
	authBL bll.AuthBL,
) *AuthHandler {
	return &AuthHandler{oauthClient: oauthClient,
		jwtService: jwtService,
		authBL:     authBL,
	}
}

func (h *AuthHandler) Mount(route *gin.Engine) {
	handler := route.Group("/auth")
	handler.GET("/google/login", server.HandlerErrorFunc(h.SignIn))
	handler.GET("/google/signup", server.HandlerErrorFunc(h.SignUp))
	handler.GET("/google/callback", server.HandlerErrorFunc(h.Callback))
}

// SignIn login using google auth
// @Summary login using google auth
// @Description login using google auth
// @Tags Auth
// @ID auth_login
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Router /auth/google/login [get]
func (h *AuthHandler) SignIn(c *gin.Context) error {

	state := base64.URLEncoding.EncodeToString([]byte(oauth.LoginState))
	conf, err := h.oauthClient.Login(c.Request.Context(), state)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, conf.URL)
	return nil
}

func (h *AuthHandler) Callback(c *gin.Context) error {
	code := c.Query(oauth.CodeNameGoogle)
	state := c.Query(oauth.StateNameGoogle)
	action, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return utils.InvalidInput.Wrap(err, "wrong state")
	}
	response, err := h.oauthClient.Callback(c.Request.Context(), code)
	if err != nil {
		return err
	}
	var user *model.User
	switch string(action) {
	case oauth.LoginState:
		user, err = h.authBL.Login(c.Request.Context(), response.Email)
	case oauth.SignUpState:
		user, err = h.authBL.SignUp(c.Request.Context(), response)
	default:
		err = fmt.Errorf("unknown state %s ", string(action))
	}
	if err != nil {
		return err
	}
	claims := &security.Identity{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		User:         user,
	}
	token, err := h.jwtService.Create(claims)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, token)
}

// SignUp register new user and tenant using google auth
// @Summary register new user and tenant using google auth
// @Description register new user and tenant using google auth
// @Tags Auth
// @ID auth_login
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Router /auth/google/login [get]
func (h *AuthHandler) SignUp(c *gin.Context) error {
	state := base64.URLEncoding.EncodeToString([]byte(oauth.SignUpState))
	conf, err := h.oauthClient.Login(c.Request.Context(), state)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, conf.URL)
	return nil
}
