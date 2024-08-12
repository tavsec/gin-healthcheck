package checks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCheck struct {
	client *redis.Client
}

func NewRedisCheck(client *redis.Client) *RedisCheck {
	return &RedisCheck{client: client}
}

func (r *RedisCheck) Pass() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.client.Ping(ctx).Result()
	return err == nil
}

func (r *RedisCheck) Name() string {
	return "redis"
}
