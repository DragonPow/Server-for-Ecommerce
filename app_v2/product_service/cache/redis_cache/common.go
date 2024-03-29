package redis_cache

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/library/math"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// parseKey:
//
//	 type: must be "product", "user",...
//	 key: must be ID or something else
//
//	Return storeKey with format "{type}/{key}". Ex: "product/1445"
func parseKey(t cache.TypeCache, k any) (storeKey string) {
	return fmt.Sprintf("%s/%v", t, k)
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

func funcConvertModel2Cache(id int64, v cache.ModelValue) (string, any) {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Marshal model2Cache fail: %v", err.Error()))
	}
	return parseKey(v.GetType(), id), data
}

func GetOne[T cache.ModelValue](r *redisCache, id int64) (T, bool) {
	rs, ok := r.base.Get(context.Background(), funcConvertId2Key[T](id))
	if !ok {
		return *new(T), false
	}
	_, v := funcConvertCache2Model[T](rs)
	return v, true
}

func GetList[T cache.ModelValue](r *redisCache, ids []int64) (map[int64]T, []int64) {
	results, err := r.base.GetList(context.Background(), math.Convert(ids, funcConvertId2Key[T]))
	if err != nil {
		return nil, ids
	}

	list := make([]any, 0, len(ids))
	miss := make([]int64, 0, len(ids))

	for idx, id := range ids {
		if results[idx] == nil {
			miss = append(miss, id)
			continue
		}
		list = append(list, results[idx])
	}

	return math.ToMap(list, funcConvertCache2Model[T]), miss
}
