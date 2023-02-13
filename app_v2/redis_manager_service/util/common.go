package util

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/library/cache/redis"
	"Server-for-Ecommerce/library/slice"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"strconv"
)

func parseKey(t cache.TypeCache, k any) (storeKey string) {
	return fmt.Sprintf("%s/%v", t, k)
}

func FuncConvertModel2Cache[T cache.ModelValue](v T) (string, any) {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Marshal model2Cache fail: %v", err.Error()))
	}
	return parseKey(v.GetType(), v.GetId()), data
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

func GetMultiple[T cache.ModelValue](r *redis.Redis, ids []int64) (values map[int64]T, missingIds []int64) {
	rs, err := r.GetList(context.Background(), slice.Map(ids, funcConvertId2Key[T]))
	if err != nil {
		return nil, ids
	}
	values = slice.KeyBy(rs, funcConvertCache2Model[T])
	missingIds = slice.Diff(ids, maps.Keys(values), func(v int64) int64 { return v })
	return values, missingIds
}
