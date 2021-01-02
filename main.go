package main

import (
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-name/config"
	"project-name/pkg/xlogger"
	"project-name/pkg/xrecover"
	"project-name/router"
	"runtime"
	"time"
)

var env string

func initArgs() {
	flag.StringVar(&env, "env", "debug", "环境变量")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err     error
		gEngine *gin.Engine
		addr    string
	)


	initArgs()
	initEnv()

	if err = config.InitConfig(env); err != nil {
		panic(err)
	}
	xlogger.InitLogger()
	defer xrecover.XRecover(nil)

	//xredis.InitRedis()
	//
	//db.InitAdapter()

	// gin init
	gin.SetMode(config.G_config.Env)
	gEngine = gin.New()
	router.InitRoute(gEngine)

	// ping 服务器
	go func() {
		if err := pingServer(); err != nil {
			xlogger.NormalLogger.Fatalf("服务器ping接口无响应: %s", err)
			return
		}

		xlogger.NormalLogger.Info("服务器ping接口成功响应")
	}()

	addr = config.G_config.Addr
	xlogger.NormalLogger.Infof("服务器监听地址: %s", addr)
	if err = http.ListenAndServe(addr, gEngine); err != nil {
		panic(err)
	}
}

func pingServer() error {
	for i := 0; i < config.G_config.MaxPingCount; i++ {
		resp, err := http.Get(config.G_config.Url + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		xlogger.NormalLogger.Info("等待服务器就绪，将在一秒后重试")
		time.Sleep(time.Second)
	}
	return errors.New("未能连接到服务器")
}
