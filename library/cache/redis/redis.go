package redis

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	client           *redis.Client
	expirationSecond time.Duration
	log              logr.Logger
}

func New(addr string, password string, expiration uint32, log logr.Logger) cache.Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	return &redisCache{
		client:           rdb,
		expirationSecond: time.Duration(expiration) * time.Second,
		log:              log,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) (string, bool) {
	result := c.client.Get(ctx, key).Val()
	if result == "" {
		return "", false
	}
	return result, true
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, key, value, c.expirationSecond).Err()
}

func (c *redisCache) GetList(ctx context.Context, keys []string) ([]interface{}, error) {
	results, err := c.client.MGet(ctx, keys...).Result()
	return results, err
}

func (c *redisCache) SetList(ctx context.Context, values map[string]interface{}) error {
	return c.client.MSet(ctx, values).Err()
}
