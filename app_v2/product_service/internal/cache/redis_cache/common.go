package redis_cache

import (
	"context"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
	"strconv"
)

// parseKey:
//   type: must be "product", "user",...
//   key: must be ID or something else
//
//  Return storeKey with format "{type}/{key}". Ex: "product/1445"
func parseKey(t cache.TypeCache, k any) (storeKey string) {
	return fmt.Sprintf("%s/%v", t, k)
}

func funcConvertId2Key[T cache.ModelValue](id int64) string {
	var t T
	v := strconv.FormatInt(id, cache.Base10Int)
	return parseKey(t.GetType(), v)
}

func funcConvertAny2Model[T cache.ModelValue](v any) (int64, T) {
	result := v.(T)
	return result.GetId(), result
}

func funcConvertModel2Any(id int64, v cache.ModelValue) (string, any) {
	return parseKey(v.GetType(), id), v
}

func GetOne[T cache.ModelValue](r *redisCache, id int64) (T, bool) {
	rs, ok := r.base.Get(context.Background(), funcConvertId2Key[T](id))
	if !ok {
		return *new(T), false
	}
	_, v := funcConvertAny2Model[T](rs)
	return v, true
}

func GetList[T cache.ModelValue](r *redisCache, ids []int64) (map[int64]T, []int64) {
	results, err := r.base.GetList(context.Background(), math.Convert(ids, funcConvertId2Key[T]))
	if err != nil {
		return nil, nil
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

	return math.ToMap(list, funcConvertAny2Model[T]), miss
}
