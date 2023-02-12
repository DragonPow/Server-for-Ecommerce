package util

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/library/cache/redis"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

func parseKey(t cache.TypeCache, k any) (storeKey string) {
	return fmt.Sprintf("%s/%v", t, k)
}

func FuncConvertModel2Cache(id int64, v cache.ModelValue) (string, any) {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Marshal model2Cache fail: %v", err.Error()))
	}
	return parseKey(v.GetType(), id), data
}

func funcConvertId2Key[T cache.ModelValue](id int64) string {
	var t T
	v := strconv.FormatInt(id, cache.Base10Int)
	return parseKey(t.GetType(), v)
}

func funcConvertCache2Model[T cache.ModelValue](v any) (int64, T) {
	var result T
	err := json.Unmarshal([]byte(v.(string)), &result)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal cache2Model fail: %v", err.Error()))
	}
	return result.GetId(), result
}

func GetOne[T cache.ModelValue](r *redis.Redis, id int64) (T, bool) {
	rs, ok := r.Get(context.Background(), funcConvertId2Key[T](id))
	if !ok {
		return *new(T), false
	}
	_, v := funcConvertCache2Model[T](rs)
	return v, true
}
