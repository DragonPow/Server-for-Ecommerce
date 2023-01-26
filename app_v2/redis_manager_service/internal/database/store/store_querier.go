package store

type StoreQuerier interface {
	Querier
	Transaction(txFunc func(StoreQuerier) error) (err error)
	Ping() error
	Close() error
}
