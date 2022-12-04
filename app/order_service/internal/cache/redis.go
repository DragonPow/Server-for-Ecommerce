package cache

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value interface{}) error
	GetList(ctx context.Context, keys []string) ([]interface{}, error)
	SetList(ctx context.Context, values map[string]interface{}) error
}

type cache struct {
	redis                 *redis.Client
	expirationMilliSecond uint32
	log                   logr.Logger
}

func New(addr string, password string, expiration uint32, log logr.Logger) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	return &cache{
		redis:                 rdb,
		expirationMilliSecond: expiration,
		log:                   log,
	}
}

func (c *cache) Get(ctx context.Context, key string) (string, bool) {
	result := c.redis.Get(ctx, key).Val()
	if result == "" {
		return "", false
	}
	return result, true
}

func (c *cache) Set(ctx context.Context, key string, value interface{}) error {
	code, err := Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, code, time.Duration(c.expirationMilliSecond)*time.Millisecond).Err()
}

func (c *cache) GetList(ctx context.Context, keys []string) ([]interface{}, error) {
	results, err := c.redis.MGet(ctx, keys...).Result()
	return results, err
}

func (c *cache) SetList(ctx context.Context, values map[string]interface{}) error {
	return c.redis.MSet(ctx, values).Err()
}
