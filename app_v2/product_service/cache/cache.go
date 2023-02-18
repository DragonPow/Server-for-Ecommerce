package cache

type Cache interface {
	GetListProduct(ids []int64) (list map[int64]Product, missIds []int64)
	GetListUser(ids []int64) (list map[int64]User, missIds []int64)
	GetListCategory(ids []int64) (list map[int64]Category, missIds []int64)
	GetListProductTemplate(ids []int64) (list map[int64]ProductTemplate, missIds []int64)
	GetListSeller(ids []int64) (list map[int64]Seller, missIds []int64)
	GetListUom(ids []int64) (list map[int64]Uom, missIds []int64)

	GetProduct(id int64) (value Product, ok bool)
	GetUser(id int64) (value User, ok bool)
	GetCategory(id int64) (value Category, ok bool)
	GetProductTemplate(id int64) (value ProductTemplate, ok bool)
	GetSeller(id int64) (value Seller, ok bool)
	GetUom(id int64) (value Uom, ok bool)

	SetMultiple(objects map[int64]ModelValue) error
	Delete(typeCache TypeCache, ids []int64) error
}

type NullCache struct {
}

func (n *NullCache) GetListProduct(ids []int64) (list map[int64]Product, missIds []int64) {
	return nil, ids
}

func (n *NullCache) GetListUser(ids []int64) (list map[int64]User, missIds []int64) {
	return nil, ids
}

func (n *NullCache) GetListCategory(ids []int64) (list map[int64]Category, missIds []int64) {
	return nil, ids
}

func (n *NullCache) GetListProductTemplate(ids []int64) (list map[int64]ProductTemplate, missIds []int64) {
	return nil, ids
}

func (n *NullCache) GetListSeller(ids []int64) (list map[int64]Seller, missIds []int64) {
	return nil, ids
}

func (n *NullCache) GetListUom(ids []int64) (list map[int64]Uom, missIds []int64) {
	return nil, ids
}

func (n *NullCache) GetProduct(id int64) (value Product, ok bool) {
	return getNullOne[Product]()
}

func (n *NullCache) GetUser(id int64) (value User, ok bool) {
	return getNullOne[User]()
}

func (n *NullCache) GetCategory(id int64) (value Category, ok bool) {
	return getNullOne[Category]()
}

func (n *NullCache) GetProductTemplate(id int64) (value ProductTemplate, ok bool) {
	return getNullOne[ProductTemplate]()
}

func (n *NullCache) GetSeller(id int64) (value Seller, ok bool) {
	return getNullOne[Seller]()
}

func (n *NullCache) GetUom(id int64) (value Uom, ok bool) {
	return getNullOne[Uom]()
}

func (n *NullCache) SetMultiple(objects map[int64]ModelValue) error {
	return nil
}

func (n *NullCache) Delete(typeCache TypeCache, ids []int64) error {
	return nil
}

func getNullOne[T ModelValue]() (v T, ok bool) {
	return *new(T), false
}
