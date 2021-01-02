package config

import (
	"github.com/spf13/viper"
)

type ConfigJwt struct {
	Secret string
	Leeway string
	Exp    string
}

type ConfigLog struct {
	Level      string //
	MaxSize    int    // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int    // 日志文件最多保存多少个备份
	MaxAge     int    // 文件最多保存多少天
	Compress   bool   // 是否压缩
}

type ConfigDb struct {
	Name           string
	Addr           string
	Username       string
	Pass           string
	DataSourceName string
	SlowLogTime    string
}

type ConfigDockerDb struct {
	Name           string
	Addr           string
	Username       string
	Pass           string
	DataSourceName string
	SlowLogTime    int
}

type ConfigRedis struct {
	Addr string
	Pass string
	Db   int
}

type Config struct {
	Env          string
	Addr         string
	Name         string
	Url          string
	MaxPingCount int
	Jwt          ConfigJwt
	Log          ConfigLog
	Db           ConfigDb
	DockerDb     ConfigDockerDb
	Redis        ConfigRedis
}

var G_config *Config

func InitConfig(env string) (err error) {
	var (
		conf Config
	)

	viper.AddConfigPath("conf")
	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&conf); err != nil {
		return
	}

	G_config = &conf
	return
}
