package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	Client *redis.Client
}

// NewRedisClient initialize the redis client
func NewRedisClient(addr, password string, db int) (*CacheClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &CacheClient{Client: client}, nil
}
