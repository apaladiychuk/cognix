package handler

import (
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler  handles authentication endpoints
type AuthHandler struct {
	oauthClient oauth.Proxy
	jwtService  security.JWTService
	userRepo    repository.UserRepository
}

func NewAuthHandler(oauthClient oauth.Proxy,
	jwtService security.JWTService,
	userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{oauthClient: oauthClient,
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

func (h *AuthHandler) Mount(route *gin.Engine) {
	handler := route.Group("/auth")
	handler.GET("/google/login", server.HandlerErrorFunc(h.SignIn))
	handler.GET("/google/callback", server.HandlerErrorFunc(h.Callback))
}

func (h *AuthHandler) SignIn(c *gin.Context) error {

	state := base64.URLEncoding.EncodeToString([]byte{})
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
	_ = state
	response, err := h.oauthClient.Callback(c.Request.Context(), code)
	if err != nil {
		return err
	}
	user, err := h.userRepo.GetByUserName(c.Request.Context(), response.User.UserName)
	if err != nil {
		return err
	}

	return server.JsonResult(c, http.StatusOK, user)
}

func (h *AuthHandler) SignUp(c *gin.Context) error {
	return nil
}
