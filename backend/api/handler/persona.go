package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PersonaHandler struct {
	personaBL bll.PersonaBL
}

func NewPersonaHandler(personaBL bll.PersonaBL) *PersonaHandler {
	return &PersonaHandler{personaBL: personaBL}
}
func (h *PersonaHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/manage/persona")
	handler.Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFunc(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFunc(h.GetByID))
	handler.POST("/", server.HandlerErrorFunc(h.Create))
	handler.PUT("/:id", server.HandlerErrorFunc(h.Update))
}

// GetAll return list of allowed personas
// @Summary return list of allowed personas
// @Description return list of allowed personas
// @Tags Persona
// @ID personas_get_all
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.Persona
// @Router /manage/persona [get]
func (h *PersonaHandler) GetAll(c *gin.Context) error {
	identity, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	personas, err := h.personaBL.GetAll(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, personas)

}

func (h *PersonaHandler) GetByID(c *gin.Context) error {
	return server.JsonResult(c, http.StatusOK, "not implemented yet")
}

func (h *PersonaHandler) Create(c *gin.Context) error {
	return server.JsonResult(c, http.StatusOK, "not implemented yet")
}

func (h *PersonaHandler) Update(c *gin.Context) error {
	return server.JsonResult(c, http.StatusOK, "not implemented yet")
}
