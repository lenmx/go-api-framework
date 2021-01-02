package token

import (
	"fmt"
	"time"
	"project-name/config"
	"project-name/pkg/errno"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// jwt token 解码后的用户上下文信息
type Context struct {
	ID       uint64
	Username string
	Nbf      int64
	Iat      int64
	Exp      int64
}

func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (data interface{}, err error) {
		var ok bool

		// Make sure the `alg` is what we except.
		if _, ok = token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = jwt.ErrSignatureInvalid
			return
		}

		data = []byte(secret)
		return
	}
}

// 验证token字符串是否有效并返回用户上下文信息
func Parse(tokenString string, secret string) (context *Context, err error) {
	var (
		token  *jwt.Token
		claims jwt.MapClaims
		ok     bool
	)
	context = &Context{}

	if token, err = jwt.Parse(tokenString, secretFunc(secret)); err != nil {
		return
	}

	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return
	}

	context.ID = uint64(claims["id"].(float64))
	context.Username = claims["username"].(string)
	context.Nbf = int64(claims["nbf"].(float64))
	context.Iat = int64(claims["iat"].(float64))
	context.Exp = int64(claims["exp"].(float64))
	return
}

// 验证gin上下文中是否有token并返回用户上下文信息
func ParseRequest(c *gin.Context) (context *Context, err error) {
	var (
		header   string
		secret   string
		tokenStr string
	)

	header = c.Request.Header.Get("Authorization")
	secret = config.G_config.Jwt.Secret
	if len(header) == 0 {
		err = errno.ErrToken
		return
	}

	fmt.Sscanf(header, "Bearer %s", &tokenStr)
	context, err = Parse(tokenStr, secret)
	return
}

// 生成token
func Sign(ctx *gin.Context, c Context, secret string) (tokenString string, err error) {
	var (
		claims jwt.MapClaims
		token  *jwt.Token
	)
	if secret == "" {
		secret = config.G_config.Jwt.Secret
	}

	_exp, err := time.ParseDuration(config.G_config.Jwt.Exp)
	if err != nil {
		panic(err)
	}

	claims = jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"nbf":      time.Now().Unix(),           // 时间不能在此之前
		"iat":      time.Now().Unix(),           // token 签发时间
		"exp":      time.Now().Add(_exp).Unix(), // 过期时间
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secret))
	return
}
