package redisactor

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"

	"github.com/redis/go-redis/v9"
)

func RedisConn(cfg *config.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb
}
