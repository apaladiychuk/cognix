package handler

import (
	"cognix.ch/api/v2/core/logic"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"sync"
)

// DocumentHandler is a type that handles document-related operations in the system.
type DocumentHandler struct {
	documentBL logic.DocumentBL
}

// NewDocumentHandler creates a new instance of DocumentHandler by injecting a DocumentBL implementation.
//
// Parameters:
// - documentBL: An instance of DocumentBL that handles document-related operations.
//
// Returns:
// - *DocumentHandler: A pointer to the newly created DocumentHandler instance.
func NewDocumentHandler(documentBL logic.DocumentBL) *DocumentHandler {
	return &DocumentHandler{documentBL: documentBL}
}

// Mount sets up the document-related routes by adding them to the specified gin.Engine instance.
func (h *DocumentHandler) Mount(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := router.Group("/api/manage/documents").Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFuncAuth(h.GetAll))
	handler.POST("/upload", server.HandlerErrorFuncAuth(h.Upload))
}

func (h *DocumentHandler) GetAll(c *gin.Context, identity *security.Identity) error {
	return nil
}

func (h *DocumentHandler) Upload(c *gin.Context, identity *security.Identity) error {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	wg := sync.WaitGroup{}
	result := make([]*parameters.DocumentUploadResponse, len(files))
	for i, f := range files {
		wg.Add(1)
		go func(idx int, file multipart.FileHeader) {
			defer wg.Done()
			fileReader, err := file.Open()
			contentType := file.Header.Get("Content-Type")
			if err != nil {
				result[idx] = &parameters.DocumentUploadResponse{
					FileName: file.Filename,
					Error:    fmt.Sprintf("open file : %s", err.Error()),
					Document: nil,
				}
				return
			}
			defer fileReader.Close()
			document, err := h.documentBL.UploadDocument(c.Request.Context(), identity.User, file.Filename, contentType, fileReader)
			if err != nil {
				result[idx] = &parameters.DocumentUploadResponse{
					FileName: file.Filename,
					Error:    fmt.Sprintf("upload document : %s", err.Error()),
					Document: nil,
				}
				return
			}
			result[idx] = &parameters.DocumentUploadResponse{
				FileName: file.Filename,
				Error:    "",
				Document: document,
			}
		}(i, *f)
	}
	wg.Wait()
	return server.JsonResult(c, http.StatusOK, result)
}
