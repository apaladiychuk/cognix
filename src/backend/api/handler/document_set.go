package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DocumentSetHandler struct {
	documentSetBL bll.DocumentSetBL
}

func NewDocumentSetHandler(documentSetBL bll.DocumentSetBL) *DocumentSetHandler {
	return &DocumentSetHandler{documentSetBL: documentSetBL}
}

func (h *DocumentSetHandler) Mount(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := router.Group("/api/manage/document_sets").Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFuncAuth(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFuncAuth(h.GetByID))
	handler.POST("/", server.HandlerErrorFuncAuth(h.Create))
	handler.PUT("/:id", server.HandlerErrorFuncAuth(h.Update))
	handler.POST("/:id/:action", server.HandlerErrorFuncAuth(h.Delete))
	handler.POST("/add_connectors", server.HandlerErrorFuncAuth(h.AddConnectors))
	handler.POST("/remove_connectors", server.HandlerErrorFuncAuth(h.RemoveConnectors))
}

// GetAll return list of document sets for current user
// @Summary return list of document sets for current user
// @Description return list of document sets for current user
// @Tags Document Set
// @ID document_set_get_by_user
// @Param archived query bool false "true for include deleted document sets"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.DocumentSet
// @Router /manage/document_sets [get]
func (h *DocumentSetHandler) GetAll(c *gin.Context, identity *security.Identity) error {
	var param parameters.ArchivedParam
	if err := c.ShouldBindQuery(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "failed to bind query params")
	}
	documentSets, err := h.documentSetBL.GetByUser(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, documentSets)
}

// GetByID return document set with connector list by id
// @Summary return document set with connector list by id
// @Description return document set with connector list by id
// @Tags Document Set
// @ID document_set_get_by_id
// @Param id path int true "document set id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.DocumentSet
// @Router /manage/document_sets/{id} [get]
func (h *DocumentSetHandler) GetByID(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	documentSet, err := h.documentSetBL.GetByID(c, identity.User, id)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, documentSet)

}

// Create creates document sets
// @Summary creates document sets
// @Description creates document sets
// @Tags Document Set
// @ID document_set_create
// @Param payload body parameters.DocumentSetParam true "document set parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.DocumentSet
// @Router /manage/document_sets [post]
func (h *DocumentSetHandler) Create(c *gin.Context, identity *security.Identity) error {
	var param parameters.DocumentSetParam
	if err := c.ShouldBindJSON(&param); err != nil {
		return utils.ErrorBadRequest.New("invalid params")
	}
	documentSet, err := h.documentSetBL.Create(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, documentSet)
}

// Update updates document sets
// @Summary updates document sets
// @Description updates document sets
// @Tags Document Set
// @ID document_set_update
// @Param id path int true "document set id"
// @Param payload body parameters.DocumentSetParam true "document set parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.DocumentSet
// @Router /manage/document_sets/{id} [put]
func (h *DocumentSetHandler) Update(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}

	var param parameters.DocumentSetParam
	if err = c.ShouldBindJSON(&param); err != nil {
		return utils.ErrorBadRequest.New("invalid params")
	}
	documentSet, err := h.documentSetBL.Update(c.Request.Context(), identity.User, id, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, documentSet)
}

// Delete delete or restore document sets
// @Summary delete or restore document sets
// @Description delete or restore document sets
// @Tags Document Set
// @ID document_set_delete
// @Param id path int true "document set id"
// @Param action path string true "action : delete | restore "
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.DocumentSet
// @Router /manage/document_sets/{id}/{action} [post]
func (h *DocumentSetHandler) Delete(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	action := c.Param("action")
	if !(action == ActionRestore || action == ActionDelete) {
		return utils.ErrorBadRequest.Newf("invalid action: should be %s or %s", ActionRestore, ActionDelete)
	}
	var documentSet *model.DocumentSet
	switch action {
	case ActionRestore:
		documentSet, err = h.documentSetBL.Restore(c.Request.Context(), identity.User, id)
	case ActionDelete:
		documentSet, err = h.documentSetBL.Delete(c.Request.Context(), identity.User, id)
	}
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, documentSet)
}

func (h *DocumentSetHandler) AddConnectors(c *gin.Context, identity *security.Identity) error {
	return nil
}

func (h *DocumentSetHandler) RemoveConnectors(c *gin.Context, identity *security.Identity) error {
	return nil
}
