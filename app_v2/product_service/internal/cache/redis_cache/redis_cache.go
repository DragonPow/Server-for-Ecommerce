package redis_cache

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache/redis"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
)

type RedisCache interface {
	cache.Cache
	Close() error
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
	return r.base.SetList(context.Background(), math.ConvertMap(objects, funcConvertModel2Any))
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
