package redis_cache

import (
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache/redis"
	"github.com/go-logr/logr"
)

type RedisCache interface {
	cache.Cache
}

type redisCache struct {
	base *redis.Redis
}

func NewCache(addr string, password string, expiration uint32, log logr.Logger) cache.Cache {
	base := redis.New(addr, password, expiration, log)
	return &redisCache{
		base: base,
	}
}

func (r *redisCache) GetListProduct(ids []int64) (list map[int64]cache.Product, missIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetListUser(ids []int64) (list map[int64]cache.User, missIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetListCategory(ids []int64) (list map[int64]cache.Category, missIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetListProductTemplate(ids []int64) (list map[int64]cache.ProductTemplate, missIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetListSeller(ids []int64) (list map[int64]cache.Seller, missIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetListUom(ids []int64) (list map[int64]cache.Uom, missIds []int64) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetProduct(id int64) (value cache.Product, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetUser(id int64) (value cache.User, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetCategory(id int64) (value cache.Category, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetProductTemplate(id int64) (value cache.ProductTemplate, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetSeller(id int64) (value cache.Seller, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (r *redisCache) GetUom(id int64) (value cache.Uom, ok bool) {
	//TODO implement me
	panic("implement me")
}
