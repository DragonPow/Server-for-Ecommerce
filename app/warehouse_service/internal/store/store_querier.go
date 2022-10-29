package store

import "context"

type StoreQuerier interface {
	Querier
	Transaction(txFunc func(StoreQuerier) error) (err error)
	Ping() error
	Close() error
	// Import
	GetImportDataBill(ctx context.Context, params GetImportDataBillParams) (*GetImportDataBillResponse, error)
}
