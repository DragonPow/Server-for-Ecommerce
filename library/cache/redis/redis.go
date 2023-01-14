package redis

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	cache.Cache
	client           *redis.Client
	expirationSecond time.Duration
	log              logr.Logger
}

func New(addr string, password string, expiration uint32, log logr.Logger) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	return &Redis{
		client:           rdb,
		expirationSecond: time.Duration(expiration) * time.Second,
		log:              log,
	}
}

func (c *Redis) Get(ctx context.Context, key string) (string, bool) {
	result := c.client.Get(ctx, key).Val()
	if result == "" {
		return "", false
	}
	return result, true
}

func (c *Redis) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, key, value, c.expirationSecond).Err()
}

func (c *Redis) GetList(ctx context.Context, keys []string) ([]interface{}, error) {
	results, err := c.client.MGet(ctx, keys...).Result()
	return results, err
}

func (c *Redis) SetList(ctx context.Context, values map[string]interface{}) error {
	return c.client.MSet(ctx, values).Err()
}
