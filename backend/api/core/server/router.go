package server

import (
	"cognix.ch/api/v2/core/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error

type JsonResponse struct {
	Status int         `json:"status,omitempty"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func HandlerErrorFunc(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			ew, ok := err.(utils.Errors)
			if !ok {
				ew.Original = err
				ew.Code = http.StatusInternalServerError
				ew.Message = err.Error()
			}
			utils.Logger.Error(ew.Original.Error())
			c.JSON(int(ew.Code), JsonResponse{
				Status: int(ew.Code),
				Error:  ew.Message,
			})

		}
	}
}

func JsonResult(c *gin.Context, status int, data interface{}) error {
	c.JSON(status, JsonResponse{
		Status: status,
		Error:  "",
		Data:   data,
	})
	return nil
}
