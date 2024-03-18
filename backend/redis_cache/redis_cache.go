package redis_cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func Set(key, value string, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}
