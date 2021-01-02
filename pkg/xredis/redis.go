package xredis

import (
	"github.com/go-redis/redis/v8"
	"project-name/config"
)

type XRedis struct {
	*redis.Client
}

var Client *XRedis

func InitRedis()  {
	var (
		options *redis.Options
		redisClient *redis.Client
	)

	options = &redis.Options{
		Addr:               config.G_config.Redis.Addr,
		Password:           config.G_config.Redis.Pass,
		DB:                 config.G_config.Redis.Db,
	}

	redisClient = redis.NewClient(options)
	Client = &XRedis{redisClient}
}