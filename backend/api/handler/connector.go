package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ConnectorHandler struct {
	connectorBL bll.ConnectorBL
}

func NewCollectorHandler(connectorRepo repository.ConnectorRepository,
	connectorBL bll.ConnectorBL) *ConnectorHandler {
	return &ConnectorHandler{
		connectorBL: connectorBL,
	}
}
func (h *ConnectorHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/api/manage/connector")
	handler.Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFunc(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFunc(h.GetById))
	handler.POST("/", server.HandlerErrorFunc(h.Create))
	handler.PUT("/:id", server.HandlerErrorFunc(h.Update))
}

// GetAll return list of allowed connectors
// @Summary return list of allowed connectors
// @Description return list of allowed connectors
// @Tags Connectors
// @ID connectors_get_all
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.Connector
// @Router /manage/connector [get]
func (h *ConnectorHandler) GetAll(c *gin.Context) error {
	identity, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	connectors, err := h.connectorBL.GetAll(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connectors)
}

// GetById return list of allowed connectors
// @Summary return list of allowed connectors
// @Description return list of allowed connectors
// @Tags Connectors
// @ID connectors_get_by_id
// @Produce  json
// @Param id path int true "connector id"
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id} [get]
func (h *ConnectorHandler) GetById(c *gin.Context) error {
	identity, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}

	connectors, err := h.connectorBL.GetByID(c.Request.Context(), identity.User, id)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connectors)
}

// Create creates connector
// @Summary creates connector
// @Description creates connector
// @Tags Connectors
// @ID connectors_create
// @Param params body parameters.CreateConnectorParam true "connector create parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.Connector
// @Router /manage/connector/ [post]
func (h *ConnectorHandler) Create(c *gin.Context) error {
	identity, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	var param parameters.CreateConnectorParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	connector, err := h.connectorBL.Create(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusCreated, connector)
}

// Update updates connector
// @Summary updates connector
// @Description updates connector
// @Tags Connectors
// @ID connectors_update
// @Param id path int true "connector id"
// @Param params body parameters.UpdateConnectorParam true "connector update parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id} [put]
func (h *ConnectorHandler) Update(c *gin.Context) error {
	identity, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	var param parameters.UpdateConnectorParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	connector, err := h.connectorBL.Update(c.Request.Context(), id, identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connector)
}
