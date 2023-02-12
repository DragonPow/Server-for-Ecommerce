package mem_cache

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"sync"
)

type myBigCache struct {
	*bigcache.BigCache
	mu              sync.RWMutex
	missMu          sync.RWMutex
	maxNumberMiss   int
	cacheMissNumber sync.Map
	maxNumberCache  int
}

func NewBigCache(maxNumberMiss, maxNumberCache int) (MemCache, error) {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         0,
		CleanWindow:        0,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		HardMaxCacheSize:   8192,
	}
	bCache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &myBigCache{
		BigCache:        bCache,
		maxNumberMiss:   maxNumberMiss,
		cacheMissNumber: sync.Map{},
		maxNumberCache:  maxNumberCache,
		mu:              sync.RWMutex{},
		missMu:          sync.RWMutex{},
	}, nil
}

func (m *myBigCache) GetListProduct(ids []int64) (values map[int64]cache.Product, missIds []int64) {
	return GetListBigCache[cache.Product](m, ids)
}

func (m *myBigCache) GetListUser(ids []int64) (values map[int64]cache.User, missIds []int64) {
	return GetListBigCache[cache.User](m, ids)
}

func (m *myBigCache) GetListCategory(ids []int64) (values map[int64]cache.Category, missIds []int64) {
	return GetListBigCache[cache.Category](m, ids)
}

func (m *myBigCache) GetListProductTemplate(ids []int64) (values map[int64]cache.ProductTemplate, missIds []int64) {
	return GetListBigCache[cache.ProductTemplate](m, ids)
}

func (m *myBigCache) GetListSeller(ids []int64) (values map[int64]cache.Seller, missIds []int64) {
	return GetListBigCache[cache.Seller](m, ids)
}

func (m *myBigCache) GetListUom(ids []int64) (values map[int64]cache.Uom, missIds []int64) {
	return GetListBigCache[cache.Uom](m, ids)
}

func (m *myBigCache) GetProduct(id int64) (value cache.Product, ok bool) {
	return GetOneBigCache[cache.Product](m, id)
}

func (m *myBigCache) GetUser(id int64) (value cache.User, ok bool) {
	return GetOneBigCache[cache.User](m, id)
}

func (m *myBigCache) GetCategory(id int64) (value cache.Category, ok bool) {
	return GetOneBigCache[cache.Category](m, id)
}

func (m *myBigCache) GetProductTemplate(id int64) (value cache.ProductTemplate, ok bool) {
	return GetOneBigCache[cache.ProductTemplate](m, id)
}

func (m *myBigCache) GetSeller(id int64) (value cache.Seller, ok bool) {
	return GetOneBigCache[cache.Seller](m, id)
}

func (m *myBigCache) GetUom(id int64) (value cache.Uom, ok bool) {
	return GetOneBigCache[cache.Uom](m, id)
}

func (m *myBigCache) SetMultiple(objects map[int64]cache.ModelValue) error {
	lockOk := m.mu.TryLock()
	if lockOk {
		defer m.mu.Unlock()
	}
	lockMissOk := m.missMu.TryLock()
	if lockMissOk {
		defer m.missMu.Unlock()
	}

	for id, value := range objects {
		key := funcConvertId2Any(id)
		storeKey := parseKey(value.GetType(), key)
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		err = m.BigCache.Set(storeKey, data)
		if err != nil {
			return err
		}
		m.cacheMissNumber.Delete(storeKey)
	}
	return nil
}

func (m *myBigCache) Delete(typeCache cache.TypeCache, ids []int64) error {
	for _, id := range ids {
		key := funcConvertId2Any(id)
		storeKey := parseKey(typeCache, key)
		m.BigCache.Delete(storeKey)
	}
	return nil
}

func (m *myBigCache) CheckAndSet(objects map[int64]cache.ModelValue) (isSet bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.missMu.Lock()
	defer m.missMu.Unlock()

	objectNeedSets := make(map[int64]cache.ModelValue, len(objects))
	for k, v := range objects {
		if m.IsMaxMiss(parseKey(v.GetType(), k)) {
			if !isSet {
				isSet = true
			}
			objectNeedSets[k] = v
		}
	}
	err = m.SetMultiple(objectNeedSets)
	return isSet, err
}

func (m *myBigCache) SetProductByAttr(object cache.Product, attrs []byte, version string) error {
	product, ok := GetOneBigCache[cache.Product](m, object.GetId())
	if !ok {
		// if not in mem cache, ignore
		return nil
	}

	err := json.Unmarshal(attrs, &product)
	if err != nil {
		return err
	}
	err = product.UpdateVersion(version)
	if err != nil {
		return err
	}
	return m.SetMultiple(map[int64]cache.ModelValue{product.GetId(): product})
}

func GetOneBigCache[T cache.ModelValue](m *myBigCache, id int64) (value T, ok bool) {
	value = *new(T)
	key := funcConvertId2Any(id)
	storeKey := parseKey(value.GetType(), key)
	data, err := m.BigCache.Get(storeKey)
	if err != nil {
		return value, false
	}
	err = json.Unmarshal(data, &value)
	if err != nil {
		panic(err)
	}
	return value, true
}

func GetListBigCache[T cache.ModelValue](m *myBigCache, ids []int64) (values map[int64]T, missingKeys []int64) {
	values = make(map[int64]T, len(ids))
	missingKeys = make([]int64, 0, len(ids))

	for _, id := range ids {
		v, ok := GetOneBigCache[T](m, id)
		if !ok {
			missingKeys = append(missingKeys, id)
			continue
		}
		values[id] = v
	}
	return
}

func (m *myBigCache) IsMaxMiss(storeKey string) bool {
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
