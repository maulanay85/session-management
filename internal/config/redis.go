package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitializeRedis(conf *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		Password: conf.RedisPassword,
		Username: conf.RedisUsername,
		DB:       0,
	})
	return rdb
}
