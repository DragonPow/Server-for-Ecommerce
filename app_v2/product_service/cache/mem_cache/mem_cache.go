package mem_cache

import (
	"encoding/json"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/cache"
	"sync"
)

type MemCache interface {
	cache.Cache
	IsMaxMiss(storeKey string) bool
	CheckAndSet(objects map[int64]cache.ModelValue) (bool, error)
	SetProductByAttr(object cache.Product, attrs []byte) error
}

type memCache struct {
	sync.Map
	mu              sync.RWMutex
	missMu          sync.RWMutex
	maxNumberMiss   int
	cacheMissNumber sync.Map
	maxNumberCache  int
}

func NewCache(maxNumberMiss, maxNumberCache int) MemCache {
	return &memCache{
		Map:             sync.Map{},
		maxNumberMiss:   maxNumberMiss,
		cacheMissNumber: sync.Map{},
		maxNumberCache:  maxNumberCache,
		mu:              sync.RWMutex{},
		missMu:          sync.RWMutex{},
	}
}

func (m *memCache) SetMultiple(objects map[int64]cache.ModelValue) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
			return
		}
	}()

	lockOk := m.mu.TryLock()
	defer func() {
		if lockOk {
			m.mu.Unlock()
		}
	}()

	for id, value := range objects {
		m.Store(value.GetType(), id, value)
	}
	return nil
}

func (m *memCache) CheckAndSet(objects map[int64]cache.ModelValue) (isSet bool, err error) {
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

func (m *memCache) GetListProduct(ids []int64) (list map[int64]cache.Product, missIds []int64) {
	return getList[cache.Product](m, ids)
}

func (m *memCache) GetListUser(ids []int64) (list map[int64]cache.User, missIds []int64) {
	return getList[cache.User](m, ids)
}

func (m *memCache) GetListCategory(ids []int64) (list map[int64]cache.Category, missIds []int64) {
	return getList[cache.Category](m, ids)
}

func (m *memCache) GetListProductTemplate(ids []int64) (list map[int64]cache.ProductTemplate, missIds []int64) {
	return getList[cache.ProductTemplate](m, ids)
}

func (m *memCache) GetListSeller(ids []int64) (list map[int64]cache.Seller, missIds []int64) {
	return getList[cache.Seller](m, ids)
}

func (m *memCache) GetListUom(ids []int64) (list map[int64]cache.Uom, missIds []int64) {
	return getList[cache.Uom](m, ids)
}

func (m *memCache) GetProduct(id int64) (cache.Product, bool) {
	return GetOne[cache.Product](m, id)
}

func (m *memCache) GetUser(id int64) (value cache.User, ok bool) {
	return GetOne[cache.User](m, id)
}

func (m *memCache) GetCategory(id int64) (value cache.Category, ok bool) {
	return GetOne[cache.Category](m, id)
}

func (m *memCache) GetProductTemplate(id int64) (value cache.ProductTemplate, ok bool) {
	return GetOne[cache.ProductTemplate](m, id)
}

func (m *memCache) GetSeller(id int64) (value cache.Seller, ok bool) {
	return GetOne[cache.Seller](m, id)
}

func (m *memCache) GetUom(id int64) (value cache.Uom, ok bool) {
	return GetOne[cache.Uom](m, id)
}

func (m *memCache) SetProductByAttr(object cache.Product, attrs []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
			return
		}
	}()

	lockOk := m.mu.TryLock()
	defer func() {
		if lockOk {
			m.mu.Unlock()
		}
	}()

	product, ok := GetOne[cache.Product](m, object.GetId())
	if !ok {
		// if not in mem cache, ignore
		return nil
	}

	err = json.Unmarshal(attrs, &product)
	if err != nil {
		return err
	}

	m.Store(object.GetType(), object.GetId(), product)
	return nil
}
