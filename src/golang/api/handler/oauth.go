package handler

import (
	"cognix.ch/api/v2/core/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

// OAuthHandler  provide oauth authentication for  third part services
type OAuthHandler struct {
}

func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{}
}

func (h *OAuthHandler) Mount(route *gin.Engine) {
	handler := route.Group("/api/oauth")
	handler.GET("/:provider/auth_url", server.HandlerErrorFunc(h.GetUrl))
	//handler.GET("/google/signup", server.HandlerErrorFunc(h.SignUp))
	handler.GET("/:provider/callback", server.HandlerErrorFunc(h.Callback))
	handler.POST("/:provider/refresh_token", server.HandlerErrorFunc(h.Refresh))
}

func (h *OAuthHandler) GetUrl(c *gin.Context) error {
	return nil
}

func (h *OAuthHandler) Callback(c *gin.Context) error {
	provider := c.Param("provider")
	query := make(map[string]string)
	c.BindQuery(&query)

	return server.JsonResult(c, http.StatusOK, c.Request.URL.String()+" ----- "+provider)
}

func (h *OAuthHandler) Refresh(c *gin.Context) error {
	return nil
}
