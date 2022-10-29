package store

import (
	"database/sql"
)

type Store struct {
	*Queries
	db *sqlx.DB
}

type StoreOptional struct {
}

func NewStore(db *sqlx.DB, option *StoreOptional) *Store {
	return &Store{db: db, Queries: New(db)}
}

func (q *Store) Ping() error {
	return q.db.Ping()
}

func (q *Store) WithTx(tx *sql.Tx) *Store {
	newRepo := *q
	newRepo.Queries = q.Queries.WithTx(tx)
	return &newRepo
}

func (q *Store) Transaction(txFunc func(StoreQuerier) error) (err error) {
	// just execute txFunc if already in transaction
	if q.Queries.tx != nil {
		return txFunc(q)
	}

	tx, err := q.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			execErr := tx.Rollback()
			if execErr != nil {
				q.log.Error(execErr, "error exec rollback")
			}
			panic(p) // re-throw panic after Rollback
		case err != nil:
			execErr := tx.Rollback() // err is non-nil; don't change it
			if execErr != nil {
				q.log.Error(execErr, "error exec rollback")
			}
		default:
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	return txFunc(q.WithTx(tx))
}
