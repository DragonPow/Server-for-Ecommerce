package mem_cache

import "github.com/DragonPow/Server-for-Ecommerce/library/math"

type MemCache interface {
	IsMaxMiss(storeKey string) bool
	GetListProduct(ids []int64) (list map[int64]Product, missIds []int64)
	GetListUser(ids []int64) (list map[int64]User, missIds []int64)
	GetListCategory(ids []int64) (list map[int64]Category, missIds []int64)
	GetListProductTemplate(ids []int64) (list map[int64]ProductTemplate, missIds []int64)
	GetListSeller(ids []int64) (list map[int64]Seller, missIds []int64)
	GetListUom(ids []int64) (list map[int64]Uom, missIds []int64)
	SetMultiple(objects map[int64]ModelValue)
}

func (m *memCache) SetMultiple(objects map[int64]ModelValue) {
	m.mu.Lock()
	for key, value := range objects {
		m.Store(value.GetType(), key, value)
	}
	m.mu.Unlock()
}

func (m *memCache) GetListProduct(ids []int64) (list map[int64]Product, missIds []int64) {
	values, miss := m.LoadMultiple(TypeProduct, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[Product](values, miss)
}

func (m *memCache) GetListUser(ids []int64) (list map[int64]User, missIds []int64) {
	values, miss := m.LoadMultiple(TypeUser, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[User](values, miss)
}

func (m *memCache) GetListCategory(ids []int64) (list map[int64]Category, missIds []int64) {
	values, miss := m.LoadMultiple(TypeCategory, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[Category](values, miss)
}

func (m *memCache) GetListProductTemplate(ids []int64) (list map[int64]ProductTemplate, missIds []int64) {
	values, miss := m.LoadMultiple(TypeProductTemplate, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[ProductTemplate](values, miss)
}

func (m *memCache) GetListSeller(ids []int64) (list map[int64]Seller, missIds []int64) {
	values, miss := m.LoadMultiple(TypeSeller, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[Seller](values, miss)
}

func (m *memCache) GetListUom(ids []int64) (list map[int64]Uom, missIds []int64) {
	values, miss := m.LoadMultiple(TypeUom, math.Convert(ids, funcConvertAny))
	return ConvertMultipleResponse[Uom](values, miss)
}
