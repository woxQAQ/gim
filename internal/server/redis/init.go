package redis

import (
	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	return rds
}
