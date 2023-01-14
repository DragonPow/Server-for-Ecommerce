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
func parseKey(t cache.TypeCache, k any) (storeKey string) {
	return fmt.Sprintf("%s/%v", t, k)
}

func (m *memCache) Store(typeObject cache.TypeCache, key any, value any) {
	storeKey := parseKey(typeObject, key)
	ok := m.mu.TryLock()
	missOk := m.missMu.TryLock()
	defer func() {
		if ok {
			m.mu.Unlock()
		}
		if missOk {
			m.missMu.Unlock()
		}
	}()

	m.cacheMissNumber.Delete(storeKey)
	m.Map.Store(storeKey, value)
}

func (m *memCache) Load(t cache.TypeCache, k any) (value any, ok bool) {
	storeKey := parseKey(t, k)
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
	lockOk := m.missMu.TryLock()
	defer func() {
		if lockOk {
			m.missMu.Unlock()
		}
	}()

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

func funcConvertId2Any(i int64) any {
	return i
}

func funcConvertAny2Model[T cache.ModelValue](v any) (int64, T) {
	result := v.(T)
	return result.GetId(), result
}

func funcConvertAny2Id(v any) int64 {
	return v.(int64)
}

// ConvertMultipleResponse
//  Convert values (format ModelValue), miss (format int64) to list (map[id]value, missIds)
//  Warn: If convert fail, panic
func ConvertMultipleResponse[T cache.ModelValue](values []any, miss []any) (map[int64]T, []int64) {
	list := math.ToMap(values, funcConvertAny2Model[T])
	missIds := math.Convert(miss, funcConvertAny2Id)
	return list, missIds
}

func GetOne[T cache.ModelValue](m *memCache, id int64) (T, bool) {
	var t T
	v, ok := m.Load(t.GetType(), id)
	if !ok {
		return *new(T), false
	}
	_, p := funcConvertAny2Model[T](v)
	return p, true
}

func getList[T cache.ModelValue](m *memCache, ids []int64) (map[int64]T, []int64) {
	var t T
	values, miss := m.LoadMultiple(t.GetType(), math.Convert(ids, funcConvertId2Any))
	return ConvertMultipleResponse[T](values, miss)
}
