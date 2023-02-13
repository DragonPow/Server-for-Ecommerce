package util

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
)

type UpdateCacheEventValue struct {
	Objects []cache.Product `json:"objects"`
}
