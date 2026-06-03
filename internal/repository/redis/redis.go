package redis

import (
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr, password string, db int) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisRepository{client: rdb}
}

// func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}) error {
// 	err := r.client.Set(ctx, key, value, 0).Err()
// 	if err != nil {
// 		return fmt.Errorf("failed to set key in Redis: %w", err)
// 	}
// 	return nil
// }

// func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
// 	val, err := r.client.Get(ctx, key).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			return "", nil // Key does not exist
// 		}
// 		return "", fmt.Errorf("failed to get key from Redis: %w", err)
// 	}
// 	return val, nil
// }
