package checks

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisCheck struct {
	Client *redis.Client
}

func (r *RedisCheck) Pass() bool {
	_, err := r.Client.Ping(context.Background()).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *RedisCheck) Name() string {
	return "redis"
}
