package cache

import (
	"context"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value interface{}) error
	GetList(ctx context.Context, keys []string) ([]interface{}, error)
	SetList(ctx context.Context, values map[string]interface{}) error
}
