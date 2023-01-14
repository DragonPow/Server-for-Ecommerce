package mem_cache

import (
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/cache"
	"github.com/DragonPow/Server-for-Ecommerce/library/math"
	"sync"
)

type MemCache interface {
	cache.Cache
	IsMaxMiss(storeKey string) bool
	SetMultiple(objects map[int64]cache.ModelValue)
}

type memCache struct {
	sync.Map
	mu              sync.RWMutex
	maxNumberMiss   int
	cacheMissNumber sync.Map
}

func NewCache(maxNumberMiss int) cache.Cache {
	return &memCache{
		Map:             sync.Map{},
		maxNumberMiss:   maxNumberMiss,
		cacheMissNumber: sync.Map{},
	}
}

func (m *memCache) SetMultiple(objects map[int64]cache.ModelValue) {
	m.mu.Lock()
	for key, value := range objects {
		m.Store(value.GetType(), key, value)
	}
	m.mu.Unlock()
}

func (m *memCache) GetListProduct(ids []int64) (list map[int64]cache.Product, missIds []int64) {
	values, miss := m.LoadMultiple(cache.TypeProduct, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[cache.Product](values, miss)
}

func (m *memCache) GetListUser(ids []int64) (list map[int64]cache.User, missIds []int64) {
	values, miss := m.LoadMultiple(cache.TypeUser, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[cache.User](values, miss)
}

func (m *memCache) GetListCategory(ids []int64) (list map[int64]cache.Category, missIds []int64) {
	values, miss := m.LoadMultiple(cache.TypeCategory, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[cache.Category](values, miss)
}

func (m *memCache) GetListProductTemplate(ids []int64) (list map[int64]cache.ProductTemplate, missIds []int64) {
	values, miss := m.LoadMultiple(cache.TypeProductTemplate, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[cache.ProductTemplate](values, miss)
}

func (m *memCache) GetListSeller(ids []int64) (list map[int64]cache.Seller, missIds []int64) {
	values, miss := m.LoadMultiple(cache.TypeSeller, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[cache.Seller](values, miss)
}

func (m *memCache) GetListUom(ids []int64) (list map[int64]cache.Uom, missIds []int64) {
	values, miss := m.LoadMultiple(cache.TypeUom, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[cache.Uom](values, miss)
}

func (m *memCache) GetProduct(id int64) (value cache.Product, ok bool) {
	v, ok := m.Load(cache.TypeProduct, id)
	if !ok {
		return cache.Product{}, false
	}
	_, p := funcConvertModel[cache.Product](v)
	return p, true
}

func (m *memCache) GetUser(id int64) (value cache.User, ok bool) {
	v, ok := m.Load(cache.TypeProduct, id)
	if !ok {
		return cache.User{}, false
	}
	_, p := funcConvertModel[cache.User](v)
	return p, true
}

func (m *memCache) GetCategory(id int64) (value cache.Category, ok bool) {
	v, ok := m.Load(cache.TypeProduct, id)
	if !ok {
		return cache.Category{}, false
	}
	_, p := funcConvertModel[cache.Category](v)
	return p, true
}

func (m *memCache) GetProductTemplate(id int64) (value cache.ProductTemplate, ok bool) {
	v, ok := m.Load(cache.TypeProduct, id)
	if !ok {
		return cache.ProductTemplate{}, false
	}
	_, p := funcConvertModel[cache.ProductTemplate](v)
	return p, true
}

func (m *memCache) GetSeller(id int64) (value cache.Seller, ok bool) {
	v, ok := m.Load(cache.TypeProduct, id)
	if !ok {
		return cache.Seller{}, false
	}
	_, p := funcConvertModel[cache.Seller](v)
	return p, true
}

func (m *memCache) GetUom(id int64) (value cache.Uom, ok bool) {
	v, ok := m.Load(cache.TypeProduct, id)
	if !ok {
		return cache.Uom{}, false
	}
	_, p := funcConvertModel[cache.Uom](v)
	return p, true
}
