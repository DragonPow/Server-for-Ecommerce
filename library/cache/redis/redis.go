package redis

import (
	"Server-for-Ecommerce/library/cache"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	cache.Cache
	client           *redis.Client
	expirationSecond time.Duration
}

type RedisOption = func(r *Redis)

func WithExpireTime(t time.Duration) RedisOption {
	return func(r *Redis) {
		r.expirationSecond = t
	}
}

func New(addr string, password string, expiration uint32) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	return &Redis{
		client:           rdb,
		expirationSecond: time.Duration(expiration) * time.Second,
	}
}

func (c *Redis) Close() error {
	return c.client.Close()
}

func (c *Redis) Get(ctx context.Context, key string) (string, bool) {
	result := c.client.Get(ctx, key).Val()
	if result == "" {
		return "", false
	}
	return result, true
}

func (c *Redis) Set(ctx context.Context, key string, value any, opts ...RedisOption) error {
	r := *c
	for _, opt := range opts {
		opt(&r)
	}
	return r.client.Set(ctx, key, value, c.expirationSecond).Err()
}

func (c *Redis) GetList(ctx context.Context, keys []string) ([]any, error) {
	results, err := c.client.MGet(ctx, keys...).Result()
	return results, err
}

func (c *Redis) SetList(ctx context.Context, values map[string]any) error {
	err := c.client.MSet(ctx, values).Err()
	if err != nil {
		return err
	}
	for k := range values {
		err = c.client.Expire(ctx, k, c.expirationSecond).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Redis) Delete(ctx context.Context, keys []string) error {
	return c.client.Del(ctx, keys...).Err()
}
