package util

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"database/sql"
	"time"
)

func ParseTimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func ParseUnixTimeToString(t sql.NullInt64) string {
	if !t.Valid {
		return EmptyString
	}
	return time.Unix(t.Int64, 0).Format(time.RFC3339)
}

func FuncConvertToCache[T cache.ModelValue](k int64, v T) (int64, cache.ModelValue) {
	return k, v
}
