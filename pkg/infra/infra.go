package infra

import (
	"web-chat/configs"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Infra struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewInfra(cfg configs.Config) *Infra {
	return &Infra{
		Redis: newRedis(cfg),
		DB:    newMysql(cfg),
	}
}
