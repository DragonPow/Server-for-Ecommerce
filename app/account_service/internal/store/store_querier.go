package store

import "context"

type StoreQuerier interface {
	Querier
	Transaction(txFunc func(StoreQuerier) error) (err error)
	Ping() error
	Close() error
	// Customer
	GetUsers(ctx context.Context, params GetUsersParams) (*GetUsersResponse, error)
}
