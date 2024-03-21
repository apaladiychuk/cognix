package handler

import (
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	oauthClient oauth.Proxy
}

func NewAuthHandler(oauthClient oauth.Proxy) *AuthHandler {
	return &AuthHandler{oauthClient: oauthClient}
}

func (h *AuthHandler) Mount(route *gin.Engine) {
	handler := route.Group("/auth")
	handler.GET("/google/login", server.HandlerErrorFunc(h.SignIn))
	handler.GET("/google/callback", server.HandlerErrorFunc(h.Callback))
}

func (h *AuthHandler) SignIn(c *gin.Context) error {
	conf, err := h.oauthClient.Login(c.Request.Context())
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, conf.URL)
	return nil
}

func (h *AuthHandler) Callback(c *gin.Context) error {
	code := c.Query(oauth.CodeNameGoogle)
	jwt, err := h.oauthClient.Callback(c.Request.Context(), code)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, jwt)
}
