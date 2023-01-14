package mem_cache

import (
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
)

// parseKey:
//   type: must be "product", "user",...
//   key: must be ID or something else
//
//  Return storeKey with format "{type}/{key}". Ex: "product/1445"
func (m *memCache) parseKey(t cache.TypeCache, k any) (storeKey string) {
	return fmt.Sprintf("%s/%v", t, k)
}

func (m *memCache) Store(typeObject cache.TypeCache, key any, value any) {
	storeKey := m.parseKey(typeObject, key)
	ok := m.mu.TryLock()
	m.cacheMissNumber.Delete(storeKey)
	m.Map.Store(storeKey, value)
	if ok {
		m.mu.Unlock()
	}
}

func (m *memCache) Load(t cache.TypeCache, k any) (value any, ok bool) {
	storeKey := m.parseKey(t, k)
	return m.Map.Load(storeKey)
}

func (m *memCache) LoadMultiple(t cache.TypeCache, keys []any) (values []any, missingKeys []any) {
	values = make([]any, len(keys))
	missingKeys = make([]any, len(keys))
	lockOk := m.mu.TryRLock()
	for _, key := range keys {
		v, ok := m.Load(t, key)
		if !ok {
			missingKeys = append(missingKeys, key)
		} else {
			values = append(values, v)
		}
	}
	if lockOk {
		m.mu.RUnlock()
	}
	return
}

// Setup Miss value
// ----------------------------------------------------

// IsMaxMiss check number miss of storeKey
//  if larger than max, delete and return true
//  if not exists or smaller than max, plus by 1
func (m *memCache) IsMaxMiss(storeKey string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, _ := m.cacheMissNumber.Load(storeKey)
	number, _ := v.(int)
	if number >= m.maxNumberMiss {
		//m.deleteMiss(storeKey)
		return true
	}
	number++
	m.cacheMissNumber.Store(storeKey, number)
	return false
}

func funcConvertAny(i int64) any {
	return i
}

func funcConvertModel[T cache.ModelValue](v any) (int64, T) {
	result := v.(T)
	return result.GetId(), result
}

func funcConvertId(v any) int64 {
	return v.(int64)
}

func ConvertMultipleResponse[T cache.ModelValue](values []any, miss []any) (map[int64]T, []int64) {
	list := math.ToMap(values, funcConvertModel[T])
	missIds := math.Convert(miss, funcConvertId)
	return list, missIds
}
