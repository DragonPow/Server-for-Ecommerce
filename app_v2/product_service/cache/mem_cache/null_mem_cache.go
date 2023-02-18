package mem_cache

import "Server-for-Ecommerce/app_v2/product_service/cache"

type NullMemCache struct {
	cache.NullCache
}

func (n *NullMemCache) CheckAndSet(objects map[int64]cache.ModelValue) (bool, error) {
	return true, nil
}

func (n *NullMemCache) SetProductByAttr(object cache.Product, attrs []byte, version string) error {
	return nil
}
