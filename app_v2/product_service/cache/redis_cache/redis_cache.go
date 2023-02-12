package redis_cache

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/library/cache/redis"
	"Server-for-Ecommerce/library/math"
	"Server-for-Ecommerce/library/slice"
	"context"
	"fmt"
	"time"
)

type RedisCache interface {
	cache.Cache
	Close() error
	GetPageProduct(page, pageSize int64, keyword string) (string, bool)
	SetPageProduct(page, pageSize int64, keyword string, data string, expireTime time.Duration) error
}

type redisCache struct {
	base *redis.Redis
}

func NewCache(addr string, password string, expiration uint32) RedisCache {
	base := redis.New(addr, password, expiration)
	return &redisCache{
		base: base,
	}
}

func (r *redisCache) Close() error {
	return r.base.Close()
}

func (r *redisCache) SetMultiple(objects map[int64]cache.ModelValue) error {
	return r.base.SetList(context.Background(), math.ConvertMap(objects, funcConvertModel2Cache))
}

func (r *redisCache) GetListProduct(ids []int64) (map[int64]cache.Product, []int64) {
	return GetList[cache.Product](r, ids)
}

func (r *redisCache) GetListUser(ids []int64) (list map[int64]cache.User, missIds []int64) {
	return GetList[cache.User](r, ids)
}

func (r *redisCache) GetListCategory(ids []int64) (list map[int64]cache.Category, missIds []int64) {
	return GetList[cache.Category](r, ids)
}

func (r *redisCache) GetListProductTemplate(ids []int64) (list map[int64]cache.ProductTemplate, missIds []int64) {
	return GetList[cache.ProductTemplate](r, ids)
}

func (r *redisCache) GetListSeller(ids []int64) (list map[int64]cache.Seller, missIds []int64) {
	return GetList[cache.Seller](r, ids)
}

func (r *redisCache) GetListUom(ids []int64) (list map[int64]cache.Uom, missIds []int64) {
	return GetList[cache.Uom](r, ids)
}

func (r *redisCache) GetProduct(id int64) (cache.Product, bool) {
	return GetOne[cache.Product](r, id)
}

func (r *redisCache) GetUser(id int64) (cache.User, bool) {
	return GetOne[cache.User](r, id)
}

func (r *redisCache) GetCategory(id int64) (cache.Category, bool) {
	return GetOne[cache.Category](r, id)
}

func (r *redisCache) GetProductTemplate(id int64) (cache.ProductTemplate, bool) {
	return GetOne[cache.ProductTemplate](r, id)
}

func (r *redisCache) GetSeller(id int64) (cache.Seller, bool) {
	return GetOne[cache.Seller](r, id)
}

func (r *redisCache) GetUom(id int64) (cache.Uom, bool) {
	return GetOne[cache.Uom](r, id)
}

func (r *redisCache) Delete(typeCache cache.TypeCache, ids []int64) error {
	keys := slice.Map(ids, func(i int64) string { return parseKey(typeCache, i) })
	return r.base.Delete(context.Background(), keys)
}

func (r *redisCache) GetPageProduct(page, pageSize int64, keyword string) (string, bool) {
	key := parseKey(cache.TypePageProduct, fmt.Sprintf("%d_%d_%s", page, pageSize, keyword))
	return r.base.Get(context.Background(), key)
}

func (r *redisCache) SetPageProduct(page, pageSize int64, keyword string, data string, expireTime time.Duration) error {
	key := parseKey(cache.TypePageProduct, fmt.Sprintf("%d_%d_%s", page, pageSize, keyword))
	return r.base.Set(context.Background(), key, data, redis.WithExpireTime(expireTime))
}
