package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PersonaHandler struct {
	personaBL bll.PersonaBL
}

func NewPersonaHandler(personaBL bll.PersonaBL) *PersonaHandler {
	return &PersonaHandler{personaBL: personaBL}
}
func (h *PersonaHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/api/manage/personas")
	handler.Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFuncAuth(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFuncAuth(h.GetByID))
	handler.POST("/", server.HandlerErrorFuncAuth(h.Create))
	handler.PUT("/:id", server.HandlerErrorFuncAuth(h.Update))
}

// GetAll return list of allowed personas
// @Summary return list of allowed personas
// @Description return list of allowed personas
// @Tags Persona
// @ID personas_get_all
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.Persona
// @Router /manage/personas [get]
func (h *PersonaHandler) GetAll(c *gin.Context, identity *security.Identity) error {
	personas, err := h.personaBL.GetAll(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, personas)

}

// GetByID return persona by id
// @Summary return persona by id
// @Description return persona by id
// @Tags Persona
// @ID persona_get_by_id
// @Param id path int true "persona id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Persona
// @Router /manage/personas/{id} [get]
func (h *PersonaHandler) GetByID(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	persona, err := h.personaBL.GetByID(c.Request.Context(), identity.User, id)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, persona)
}

// Create create persona
// @Summary create persona
// @Description create persona
// @Tags Persona
// @ID persona_create
// @Param id body parameters.PersonaParam true "persona payload"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.Persona
// @Router /manage/personas [post]
func (h *PersonaHandler) Create(c *gin.Context, identity *security.Identity) error {
	var param parameters.PersonaParam
	if err := c.ShouldBindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	persona, err := h.personaBL.Create(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusCreated, persona)
}

// Update update persona
// @Summary update persona
// @Description update persona
// @Tags Persona
// @ID persona_update
// @Param id path int true "persona id"
// @Param id body parameters.PersonaParam true "persona payload"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Persona
// @Router /manage/personas/{id} [put]
func (h *PersonaHandler) Update(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	var param parameters.PersonaParam
	if err = c.ShouldBindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	persona, err := h.personaBL.Update(c.Request.Context(), id, identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, persona)
}
