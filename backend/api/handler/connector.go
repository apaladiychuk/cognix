package handler

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"strconv"
	"time"
)

type ConnectorHandler struct {
	connectorRepo  repository.ConnectorRepository
	credentialRepo repository.CredentialRepository
}

func NewCollectorHandler(connectorRepo repository.ConnectorRepository,
	credentialRepo repository.CredentialRepository) *ConnectorHandler {
	return &ConnectorHandler{connectorRepo: connectorRepo,
		credentialRepo: credentialRepo,
	}
}
func (h *ConnectorHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/manage/connector")
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
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	connectors, err := h.connectorRepo.GetAll(c.Request.Context(), claims.TenantID, claims.UserID)
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
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id} [get]
func (h *ConnectorHandler) GetById(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	connectors, err := h.connectorRepo.GetAll(c.Request.Context(), claims.TenantID, claims.UserID)
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
// @Param params body CreateConnectorParam true "connector create parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.Connector
// @Router /manage/connector/ [post]
func (h *ConnectorHandler) Create(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	var param CreateConnectorParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	cred, err := h.credentialRepo.GetByID(c.Request.Context(), param.CredentialID, claims.TenantID, claims.UserID)
	if err != nil {
		return err
	}
	if cred.Source != param.Source {
		return utils.InvalidInput.New("wrong credential source")
	}
	connector := model.Connector{
		CredentialID:            param.CredentialID,
		Name:                    param.Name,
		Source:                  param.Source,
		InputType:               param.InputType,
		ConnectorSpecificConfig: param.ConnectorSpecificConfig,
		RefreshFreq:             param.RefreshFreq,
		UserID:                  uuid.MustParse(claims.UserID),
		TenantID:                uuid.MustParse(claims.TenantID),
		Shared:                  param.Shared,
		Disabled:                param.Disabled,
		CreatedDate:             time.Now().UTC(),
	}
	if err = h.connectorRepo.Create(c, &connector); err != nil {
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
// @Param params body UpdateConnectorParam true "connector update parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id} [put]
func (h *ConnectorHandler) Update(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	var param UpdateConnectorParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}

	connector, err := h.connectorRepo.GetByID(c.Request.Context(), id, claims.TenantID, claims.UserID)
	if err != nil {
		return err
	}
	cred, err := h.credentialRepo.GetByID(c.Request.Context(), param.CredentialID, claims.TenantID, claims.UserID)
	if err != nil {
		return err
	}
	if cred.Source != connector.Source {
		return utils.InvalidInput.New("wrong credential source")
	}
	connector.ConnectorSpecificConfig = param.ConnectorSpecificConfig
	connector.CredentialID = param.CredentialID
	connector.Name = param.Name
	connector.InputType = param.InputType
	connector.RefreshFreq = param.RefreshFreq
	connector.Shared = param.Shared
	connector.Disabled = param.Disabled
	connector.UpdatedDate = null.TimeFrom(time.Now().UTC())
	if err = h.connectorRepo.Update(c.Request.Context(), connector); err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connector)
}
