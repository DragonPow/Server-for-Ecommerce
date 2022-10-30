// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package store

import (
	"database/sql"
	"time"
)

type Account struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	CreateDate time.Time `json:"create_date"`
	WriteDate  time.Time `json:"write_date"`
}

type CustomerInfo struct {
	ID         int64          `json:"id"`
	AccountID  int64          `json:"account_id"`
	Name       string         `json:"name"`
	Phone      sql.NullString `json:"phone"`
	Address    sql.NullString `json:"address"`
	CreateDate time.Time      `json:"create_date"`
	WriteDate  time.Time      `json:"write_date"`
}