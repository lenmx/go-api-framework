package middleware

import (
	"time"
	"project-name/config"
	"project-name/handler"
	"project-name/pkg/errno"
	"project-name/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext, err := token.ParseRequest(c)
		if err != nil {
			handler.Response(c, errno.ErrToken, nil)
			c.Abort()
			return
		}

		if !validNbf(userContext.Nbf, config.G_config.Jwt.Leeway) {
			handler.Response(c, errno.ErrToken, nil )
			c.Abort()
			return
		}

		if !validExp(userContext.Exp) {
			handler.Response(c, errno.ErrToken, nil )
			c.Abort()
			return
		}

		c.Next()
	}
}

func validNbf(nbf int64, leeway string) (bool) {
	now := time.Now()
	leewayDuration, err := time.ParseDuration(leeway)
	if err != nil && leeway != "" {
		return false
	}

	_nbf := time.Unix(nbf, 0)
	if now.Add(leewayDuration).Before(_nbf) {
		return false
	}

	return true
}

func validExp(exp int64) (bool) {
	now := time.Now()
	_exp := time.Unix(exp, 0)

	if now.After(_exp) {
		return false
	}

	return true
}
