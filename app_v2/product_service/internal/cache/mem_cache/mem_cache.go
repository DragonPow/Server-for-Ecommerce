package mem_cache

import (
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
	"sync"
)

type MemCache interface {
	cache.Cache
	IsMaxMiss(storeKey string) bool
	CheckAndSet(objects map[int64]cache.ModelValue) (bool, error)
}

type memCache struct {
	sync.Map
	mu              sync.RWMutex
	missMu          sync.RWMutex
	maxNumberMiss   int
	cacheMissNumber sync.Map
	maxNumberCache  int
}

func (m *memCache) GetList(t cache.ModelValue, ids []int64) (list map[int64]cache.ModelValue, missIds []int64) {
	switch t.GetType() {
	case cache.TypeProduct:
		var l map[int64]cache.Product
		l, missIds = getList[cache.Product](m, ids)
		return math.ConvertMap(l, util.FuncConvertToCache[cache.Product]), missIds
	case cache.TypeProductTemplate:
		var l map[int64]cache.ProductTemplate
		l, missIds = getList[cache.ProductTemplate](m, ids)
		return math.ConvertMap(l, util.FuncConvertToCache[cache.ProductTemplate]), missIds
	case cache.TypeUom:
		var l map[int64]cache.Uom
		l, missIds = getList[cache.Uom](m, ids)
		return math.ConvertMap(l, util.FuncConvertToCache[cache.Uom]), missIds
	case cache.TypeCategory:
		var l map[int64]cache.Category
		l, missIds = getList[cache.Category](m, ids)
		return math.ConvertMap(l, util.FuncConvertToCache[cache.Category]), missIds
	case cache.TypeSeller:
		var l map[int64]cache.Seller
		l, missIds = getList[cache.Seller](m, ids)
		return math.ConvertMap(l, util.FuncConvertToCache[cache.Seller]), missIds
	case cache.TypeUser:
		var l map[int64]cache.User
		l, missIds = getList[cache.User](m, ids)
		return math.ConvertMap(l, util.FuncConvertToCache[cache.User]), missIds
	}
	return
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
