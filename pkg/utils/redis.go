package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis(config Configuration) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       0,
	})
}

func CloseRedis() error {
	return RDB.Close()
}
