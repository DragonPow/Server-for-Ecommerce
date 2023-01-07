package util

import (
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
