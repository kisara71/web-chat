package infra

import (
	"web-chat/configs"

	"github.com/redis/go-redis/v9"
)

func newRedis(cfg configs.Config) *redis.Client {
	redisConf := cfg.RedisConf
	return redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})
}
