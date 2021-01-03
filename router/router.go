package router

import (
	"net/http"
	"project-name/config"
	"project-name/handler"
	"project-name/handler/common"
	"project-name/pkg/token"

	"project-name/handler/sd"
	"project-name/router/middleware"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "project-name/docs"
)

func InitRoute() (g *gin.Engine) {
	g = gin.New()
	gin.SetMode(config.G_config.Env)

	UseMiddleware(g)
	Load(g)

	return g
}

func UseMiddleware(g *gin.Engine) {
	g.Use(middleware.ErrHandler())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestId())
}

func Load(g *gin.Engine, mw ...gin.HandlerFunc) {
	g.NoMethod(func(context *gin.Context) {
		context.String(http.StatusMethodNotAllowed, "The method not allowed.")
	})
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	commonGroup := g.Group("v1/common")
	//commonGroup.Use(middleware.AuthMiddleware())
	{
		commonGroup.GET("/ip/:ip", common.GetIpInfo)
	}

	testUserGroup := g.Group("test/user")
	{
		testUserGroup.GET("/login", func(context *gin.Context) {
			userContext := &token.Context{
				ID:       1,
				Username: "larry",
			}
			token, err := token.Sign(context, *userContext, config.G_config.Jwt.Secret)
			if err != nil {
				panic(err)
			}

			handler.Response(context, nil, token)
		})
	}

	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return
}
