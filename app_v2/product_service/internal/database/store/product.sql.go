// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: product.sql

package store

import (
	"context"

	"github.com/lib/pq"
)

const getProducts = `-- name: getProducts :many
SELECT id, template_id, name, origin_price, sale_price, state, variants, create_uid, create_date, write_uid, write_time
FROM product
WHERE
    CASE WHEN array_length($1::int8[], 1) > 0 THEN id = ANY($1::int8[]) ELSE TRUE END
AND CASE WHEN $2::int8 > 0 THEN template_id = $2::int8 ELSE TRUE END
`

type getProductsParams struct {
	Ids               []int64 `json:"ids"`
	ProductTemplateID int64   `json:"product_template_id"`
}

func (q *Queries) getProducts(ctx context.Context, arg getProductsParams) ([]Product, error) {
	rows, err := q.query(ctx, q.getProductsStmt, getProducts, pq.Array(arg.Ids), arg.ProductTemplateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.TemplateID,
			&i.Name,
			&i.OriginPrice,
			&i.SalePrice,
			&i.State,
			&i.Variants,
			&i.CreateUid,
			&i.CreateDate,
			&i.WriteUid,
			&i.WriteTime,
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
