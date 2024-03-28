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

type CredentialHandler struct {
	credentialRepo repository.CredentialRepository
}

func NewCredentialHandler(credentialRepo repository.CredentialRepository) *CredentialHandler {
	return &CredentialHandler{credentialRepo: credentialRepo}
}
func (h *CredentialHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/manage/credential")
	handler.Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFunc(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFunc(h.GetByID))
	handler.POST("/", server.HandlerErrorFunc(h.Create))
	handler.PUT("/:id", server.HandlerErrorFunc(h.Update))
}

// GetAll return list of allowed credentials
// @Summary return list of allowed credentials
// @Description return list of allowed credentials
// @Tags Credentials
// @ID credentials_get_all
// @param source query string false "source of credentials"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.Credential
// @Router /manage/credential [get]
func (h *CredentialHandler) GetAll(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	source := c.Query("source")
	credentials, err := h.credentialRepo.GetAll(c.Request.Context(), claims.TenantID, claims.UserID, source)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, credentials)
}

// GetByID return credential by id
// @Summary return list of allowed credentials
// @Description return list of allowed credentials
// @Tags Credentials
// @ID credentials_get_by_id
// @Param id path int true "credential id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Credential
// @Router /manage/credential/{id} [get]
func (h *CredentialHandler) GetByID(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}

	credential, err := h.credentialRepo.GetByID(c.Request.Context(), id, claims.TenantID, claims.UserID)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, credential)
}

// Create creates new credential
// @Summary creates new credential
// @Description creates new credential
// @Tags Credentials
// @ID credentials_create
// @Param params body CreateCredentialParam true "credential create parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.Credential
// @Router /manage/credential [post]
func (h *CredentialHandler) Create(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	var param CreateCredentialParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	credential := model.Credential{
		UserID:         uuid.MustParse(claims.UserID),
		TenantID:       uuid.MustParse(claims.TenantID),
		Source:         param.Source,
		CreatedDate:    time.Now().UTC(),
		Shared:         param.Shared,
		CredentialJson: param.CredentialJson,
	}
	if err = h.credentialRepo.Create(c.Request.Context(), &credential); err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusCreated, credential)
}

// Update updates credential
// @Summary updates credential
// @Description updates credential
// @Tags Credentials
// @ID credentials_update
// @Param id path int true "credential id"
// @Param params body UpdateCredentialParam true "credential update parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Credential
// @Router /manage/credential/{id} [put]
func (h *CredentialHandler) Update(c *gin.Context) error {
	claims, err := server.GetContextClaims(c)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	var param UpdateCredentialParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	credential, err := h.credentialRepo.GetByID(c.Request.Context(), id, claims.TenantID, claims.UserID)
	if err != nil {
		return err
	}
	if credential.UserID.String() != claims.UserID {
		return utils.ErrorPermission.New("you are not credential owner.")
	}
	credential.CredentialJson = param.CredentialJson
	credential.Shared = param.Shared
	credential.UpdatedDate = null.TimeFrom(time.Now().UTC())
	if err = h.credentialRepo.Update(c.Request.Context(), credential); err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, credential)
}
