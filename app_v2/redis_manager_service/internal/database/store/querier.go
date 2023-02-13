// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package store

import (
	"context"
)

type Querier interface {
	GetCategories(ctx context.Context, ids []int64) ([]Category, error)
	GetProductAndRelations(ctx context.Context, ids []int64) ([]GetProductAndRelationsRow, error)
	GetProductTemplates(ctx context.Context, ids []int64) ([]ProductTemplate, error)
	GetProducts(ctx context.Context, ids []int64) ([]Product, error)
	GetSellers(ctx context.Context, ids []int64) ([]Seller, error)
	GetUoms(ctx context.Context, ids []int64) ([]Uom, error)
	GetUsers(ctx context.Context, ids []int64) ([]User, error)
}

var _ Querier = (*Queries)(nil)
