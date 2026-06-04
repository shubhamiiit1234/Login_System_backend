package database

import (
	"github.com/go-redis/redis/v8"
)

func InitializeRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return rdb
}
