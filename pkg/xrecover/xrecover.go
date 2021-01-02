package xrecover

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project-name/handler"
	"project-name/pkg/errno"
	"project-name/pkg/xlogger"
)

func XRecover(context *gin.Context) {
	if data := recover(); data != nil {
		var err *errno.Errno

		switch typed := data.(type) {
		case errno.Errno:
			err = &typed
		case error:
			err = errno.ConvertErr(typed)
		default:
			err = errno.InternalServerError
		}

		xlogger.ErrorLogger.Error("", zap.NamedError("err", err))

		if context != nil {
			handler.Response(context, err, nil)
			context.Abort()
		}

		return
	}
}
