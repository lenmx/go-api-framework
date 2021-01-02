package middleware

import (
	"github.com/gin-gonic/gin"
	"project-name/pkg/xrecover"
)

func ErrHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer xrecover.XRecover(context)
		context.Next()
	}
}
