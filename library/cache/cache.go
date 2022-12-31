package cache

import (
	"context"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value any) error
	GetList(ctx context.Context, keys []string) ([]any, error)
	SetList(ctx context.Context, values map[string]any) error
}
