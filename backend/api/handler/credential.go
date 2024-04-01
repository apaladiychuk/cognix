package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CredentialHandler struct {
	credentialBl bll.CredentialBL
}

func NewCredentialHandler(credentialBl bll.CredentialBL) *CredentialHandler {
	return &CredentialHandler{credentialBl: credentialBl}
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
	ident, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	source := c.Query("source")
	credentials, err := h.credentialBl.GetAll(c.Request.Context(), ident.User, source)
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
	ident, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}

	credential, err := h.credentialBl.GetByID(c.Request.Context(), ident.User, id)
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
// @Param params body parameters.CreateCredentialParam true "credential create parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.Credential
// @Router /manage/credential [post]
func (h *CredentialHandler) Create(c *gin.Context) error {
	ident, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	var param parameters.CreateCredentialParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	credential, err := h.credentialBl.Create(c.Request.Context(), ident.User, &param)
	if err != nil {
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
// @Param params body parameters.UpdateCredentialParam true "credential update parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Credential
// @Router /manage/credential/{id} [put]
func (h *CredentialHandler) Update(c *gin.Context) error {
	ident, err := server.GetContextIdentity(c)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.InvalidInput.New("id should be presented")
	}
	var param parameters.UpdateCredentialParam
	if err = c.BindJSON(&param); err != nil {
		return utils.InvalidInput.Wrap(err, "wrong payload")
	}
	credential, err := h.credentialBl.Update(c.Request.Context(), id, ident.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, credential)
}
