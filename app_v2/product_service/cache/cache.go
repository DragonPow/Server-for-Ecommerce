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
