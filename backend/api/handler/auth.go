package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// AuthHandler  handles authentication endpoints
type AuthHandler struct {
	oauthClient oauth.Proxy
	jwtService  security.JWTService
	authBL      bll.AuthBL
	storage     storage.Storage
}

func NewAuthHandler(oauthClient oauth.Proxy,
	jwtService security.JWTService,
	authBL bll.AuthBL,
	storage storage.Storage,
) *AuthHandler {
	return &AuthHandler{oauthClient: oauthClient,
		jwtService: jwtService,
		authBL:     authBL,
		storage:    storage,
	}
}

func (h *AuthHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/auth")
	handler.GET("/google/login", server.HandlerErrorFunc(h.SignIn))
	handler.GET("/google/signup", server.HandlerErrorFunc(h.SignUp))
	handler.GET("/google/callback", server.HandlerErrorFunc(h.Callback))
	handler.GET("/google/invite", server.HandlerErrorFunc(h.JoinToTenant))
	adminHandler := route.Group("/auth")
	adminHandler.Use(authMiddleware)
	adminHandler.POST("/google/invite", server.HandlerErrorFuncAuth(h.Invite))
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
	buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.LoginState})
	if err != nil {
		return utils.Internal.Wrap(err, "can not marshal payload")
	}
	state := base64.URLEncoding.EncodeToString(buf)
	url, err := h.oauthClient.Login(c.Request.Context(), state)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, url)
	return nil
}

func (h *AuthHandler) Callback(c *gin.Context) error {
	code := c.Query(oauth.CodeNameGoogle)

	buf, err := base64.URLEncoding.DecodeString(c.Query(oauth.StateNameGoogle))
	if err != nil {
		return utils.InvalidInput.Wrap(err, "wrong state")
	}
	var state parameters.OAuthParam
	if err = json.Unmarshal(buf, &state); err != nil {
		return utils.Internal.Wrap(err, "can not unmarshal OAuth state")
	}

	response, err := h.oauthClient.Callback(c.Request.Context(), code)
	if err != nil {
		return err
	}
	var user *model.User
	switch state.Action {
	case oauth.LoginState:
		user, err = h.authBL.Login(c.Request.Context(), response.Email)
	case oauth.SignUpState:
		user, err = h.authBL.SignUp(c.Request.Context(), response)
	default:
		err = fmt.Errorf("unknown state %s ", state.Action)
	}
	if err != nil {
		return err
	}
	identity := &security.Identity{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		User:         user,
	}
	token, err := h.jwtService.Create(identity)
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
	buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.SignUpState})
	if err != nil {
		return utils.Internal.Wrap(err, "can not marshal payload")
	}

	state := base64.URLEncoding.EncodeToString(buf)
	url, err := h.oauthClient.Login(c.Request.Context(), state)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, url)
	return nil
}

func (h *AuthHandler) Invite(c *gin.Context, identity *security.Identity) error {
	var param parameters.InviteParam
	if err := c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "can not parse payload")
	}
	if err := param.Validate(); err != nil {
		return utils.InvalidInput.Wrap(err, err.Error())
	}
	buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.InviteState,
		Role:     param.Role,
		Email:    param.Email,
		TenantID: identity.User.TenantID.String(),
	})
	if err != nil {
		return utils.Internal.Wrap(err, "can not marshal payload")
	}
	key := uuid.New()
	if err = h.storage.Save(key.String(), buf); err != nil {
		return err
	}
	state := base64.URLEncoding.EncodeToString([]byte(key.String()))

	return server.JsonResult(c, http.StatusOK, fmt.Sprintf("%s?state=%s", param.BaseURL, state))
}

func (h *AuthHandler) JoinToTenant(c *gin.Context) error {
	param := c.Query("state")

	key, err := base64.URLEncoding.DecodeString(param)
	if err != nil {
		return utils.InvalidInput.Wrap(err, "wrong state")
	}
	value, err := h.storage.GetValue(string(key))

	state := base64.URLEncoding.EncodeToString(value)

	url, err := h.oauthClient.Login(c.Request.Context(), state)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, url)
	return nil
}
