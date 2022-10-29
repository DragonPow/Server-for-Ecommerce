// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package store

import (
	"database/sql"
)

type ImportBill struct {
	ID                int64          `json:"id"`
	Code              string         `json:"code"`
	Note              sql.NullString `json:"note"`
	ContactPersonName sql.NullString `json:"contact_person_name"`
	ContactEmail      sql.NullString `json:"contact_email"`
	ContactPhone      sql.NullString `json:"contact_phone"`
	CreateUid         sql.NullInt64  `json:"create_uid"`
	CreateDate        sql.NullTime   `json:"create_date"`
	WriteUid          sql.NullInt64  `json:"write_uid"`
	WriteDate         sql.NullTime   `json:"write_date"`
}

type ImportBillDetail struct {
	ID         int64        `json:"id"`
	ImportID   int64        `json:"import_id"`
	ProductID  int64        `json:"product_id"`
	Quantity   float64      `json:"quantity"`
	CreateDate sql.NullTime `json:"create_date"`
	WriteDate  sql.NullTime `json:"write_date"`
}
