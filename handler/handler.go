package handler

import (
	"net/http"

	"project-name/pkg/errno"

	"github.com/gin-gonic/gin"
)

type XResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)
	c.JSON(http.StatusOK, XResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
