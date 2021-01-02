package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
	"project-name/handler"
	"project-name/pkg/constvar"
	"project-name/pkg/errno"
	"project-name/pkg/xlogger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			start, end       time.Time
			ts               time.Duration
			path, method, ip string

			blw      *bodyLogWriter
			response handler.XResponse
			code     int
			message  string
			err      error
		)

		start = time.Now()
		path = c.Request.URL.Path
		method = c.Request.Method
		ip = c.ClientIP()

		// 读取body并还原
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		blw = &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// Continue
		c.Next()

		// 计算耗时
		end = time.Now()
		ts = end.Sub(start)

		if err = json.Unmarshal([]byte(blw.body.Bytes()), &response); err != nil {
			code = errno.InternalServerError.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		// 审计日志
		xlogger.MonitorLogger.Info(
			"",
			zap.String("start", start.Format(constvar.DateTimeFormatMs)),
			zap.String("end", end.Format(constvar.DateTimeFormatMs)),
			zap.Duration("timespan", ts),
			zap.String("path", path),
			zap.String("method", method),
			zap.String("ip", ip),
			zap.Int("code", code),
			zap.String("message", message),
		)
	}
}
