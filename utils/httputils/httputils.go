package httputils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ResponseOK  = "OK"
	ResponseErr = "Error"
)

type HttpResponse struct {
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func ResponseOk(c *gin.Context, data interface{}, desc string) {
	c.JSON(http.StatusOK, HttpResponse{
		Description: desc,
		Status:      ResponseOK,
		Data:        data,
	})
}
func ResponseError(c *gin.Context, data interface{}, desc string) {

	c.JSON(http.StatusOK, HttpResponse{
		Description: desc,
		Status:      ResponseErr,
		Data:        data,
	})
}
