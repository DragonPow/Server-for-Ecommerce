// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: import.sql

package store

import (
	"context"

	"github.com/lib/pq"
)

const getImportBillDetails = `-- name: getImportBillDetails :many
SELECT id, import_id, product_id, quantity, create_date, write_date
FROM import_bill_detail
WHERE
    CASE WHEN array_length($1::int8[],1) > 0 THEN id = ANY($1::int8[]) ELSE TRUE END
AND CASE WHEN array_length($2::int8[],1) > 0 THEN import_id = ANY($2::int8[]) ELSE TRUE END
`

type getImportBillDetailsParams struct {
	Ids      []int64 `json:"ids"`
	ImportID []int64 `json:"import_id"`
}

func (q *Queries) getImportBillDetails(ctx context.Context, arg getImportBillDetailsParams) ([]ImportBillDetail, error) {
	rows, err := q.query(ctx, q.getImportBillDetailsStmt, getImportBillDetails, pq.Array(arg.Ids), pq.Array(arg.ImportID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImportBillDetail{}
	for rows.Next() {
		var i ImportBillDetail
		if err := rows.Scan(
			&i.ID,
			&i.ImportID,
			&i.ProductID,
			&i.Quantity,
			&i.CreateDate,
			&i.WriteDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getImportBills = `-- name: getImportBills :many
SELECT id, code, note, contact_person_name, contact_email, contact_phone, create_uid, create_date, write_uid, write_date
FROM import_bill
WHERE id = ANY($1::int8[])
`

func (q *Queries) getImportBills(ctx context.Context, ids []int64) ([]ImportBill, error) {
	rows, err := q.query(ctx, q.getImportBillsStmt, getImportBills, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ImportBill{}
	for rows.Next() {
		var i ImportBill
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.Note,
			&i.ContactPersonName,
			&i.ContactEmail,
			&i.ContactPhone,
			&i.CreateUid,
			&i.CreateDate,
			&i.WriteUid,
			&i.WriteDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
