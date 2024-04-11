package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TenantHandler  handles authentication endpoints
type TenantHandler struct {
	TenantBL bll.TenantBL
}

func NewTenantHandler(TenantBL bll.TenantBL) *TenantHandler {
	return &TenantHandler{
		TenantBL: TenantBL,
	}
}

func (h *TenantHandler) Mount(route *gin.Engine, TenantMiddleware gin.HandlerFunc) {
	handler := route.Group("/tenant")
	handler.Use(TenantMiddleware)
	handler.GET("/users", server.HandlerErrorFuncAuth(h.GetUserList))
}

// GetUserList return list of users
// @Summary return list of users
// @Description return list of users
// @Tags Tenant
// @ID tenant_get_users
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.User
// @Router /tenant/users [get]
func (h *TenantHandler) GetUserList(c *gin.Context, identity *security.Identity) error {
	users, err := h.TenantBL.GetUsers(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, users)
}
