// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: seller.sql

package store

import (
	"context"
	"database/sql"
)

const createSeller = `-- name: CreateSeller :one
INSERT INTO seller
(name, description, phone, address, logo_url, manager_id, create_uid, create_date, write_uid, write_date)
VALUES ($1, $2, $3, $4, $5, $6,
        case when $7::int8 > 0 then $7::int8 else 1 end,
        case when $7::int8 > 0 then $7::int8 else 1 end,
        now() AT TIME ZONE 'utc',
        now() AT TIME ZONE 'utc') RETURNING id
`

type CreateSellerParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Phone       sql.NullString `json:"phone"`
	Address     sql.NullString `json:"address"`
	LogoUrl     sql.NullString `json:"logo_url"`
	ManagerID   int64          `json:"manager_id"`
	CreateUid   int64          `json:"create_uid"`
}

func (q *Queries) CreateSeller(ctx context.Context, arg CreateSellerParams) (int64, error) {
	row := q.queryRow(ctx, q.createSellerStmt, createSeller,
		arg.Name,
		arg.Description,
		arg.Phone,
		arg.Address,
		arg.LogoUrl,
		arg.ManagerID,
		arg.CreateUid,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateSeller = `-- name: UpdateSeller :exec
UPDATE seller
SET name        = $1,
    description = $2,
    phone       = $3,
    address     = $4,
    logo_url    = $5,
    manager_id  = $6,
    write_uid   = case when $7::int8 > 0 then $7::int8 else 1 end,
    write_date  = now() AT TIME ZONE 'utc'
WHERE id = $8::int8
`

type UpdateSellerParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Phone       sql.NullString `json:"phone"`
	Address     sql.NullString `json:"address"`
	LogoUrl     sql.NullString `json:"logo_url"`
	ManagerID   int64          `json:"manager_id"`
	CreateUid   int64          `json:"create_uid"`
	ID          int64          `json:"id"`
}

func (q *Queries) UpdateSeller(ctx context.Context, arg UpdateSellerParams) error {
	_, err := q.exec(ctx, q.updateSellerStmt, updateSeller,
		arg.Name,
		arg.Description,
		arg.Phone,
		arg.Address,
		arg.LogoUrl,
		arg.ManagerID,
		arg.CreateUid,
		arg.ID,
	)
	return err
}