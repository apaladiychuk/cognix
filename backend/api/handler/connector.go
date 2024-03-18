package handler

import (
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConnectorHandler struct {
	connectorRepo repository.ConnectorRepository
}

func NewCollectorHandler(connectorRepo repository.ConnectorRepository) *ConnectorHandler {
	return &ConnectorHandler{connectorRepo: connectorRepo}
}
func (h *ConnectorHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/admin/connector")
	handler.Use(authMiddleware)
	handler.POST("/", server.HandlerErrorFunc(h.CreateConnector))
}

func (h *ConnectorHandler) CreateConnector(c *gin.Context) error {
	//todo

	return server.JsonResult(c, http.StatusOK, "ok")
}
