package redis_cache

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/app_v2/product_service/util"
	"time"
)

type NullRedisCache struct {
	cache.NullCache
}

func (n *NullRedisCache) Close() error {
	return nil
}

func (n *NullRedisCache) GetPageProduct(page, pageSize int64, keyword string) (string, bool) {
	return util.EmptyString, false
}

func (n *NullRedisCache) SetPageProduct(page, pageSize int64, keyword string, data string, expireTime time.Duration) error {
	return nil
}

func (n *NullRedisCache) DeletePage(page, pageSize int64, keyword string) error {
	return nil
}
