package checks

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCheck struct {
	Title  string
	Client *redis.Client
}

func NewRedisCheck(client *redis.Client) RedisCheck {
	check := RedisCheck{Client: client}
	return check
}

func (r *RedisCheck) Pass() bool {
	_, err := r.Client.Ping(context.Background()).Result()
	return err == nil
}

func (r *RedisCheck) Name() string {
	if r.Title != "" {
		return r.Title
	}
	return "redis"
}
