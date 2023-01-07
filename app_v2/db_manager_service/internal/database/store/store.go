package store

import (
	"database/sql"
	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	*Queries
	db  *sqlx.DB
	log logr.Logger
}

func NewStore(db *sqlx.DB, log logr.Logger) *Store {
	return &Store{
		db:      db,
		log:     log,
		Queries: New(db),
	}
}

func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) WithTx(tx *sql.Tx) *Store {
	newRepo := *s
	//newRepo.Queries = s.Queries.WithTx(tx)
	return &newRepo
}

func (s *Store) Transaction(txFunc func(StoreQuerier) error) (err error) {
	// just execute txFunc if already in transaction
	if s.Queries.tx != nil {
		return txFunc(s)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			execErr := tx.Rollback()
			if execErr != nil {
				s.log.Error(execErr, "error exec rollback")
			}
			panic(p) // re-throw panic after Rollback
		case err != nil:
			execErr := tx.Rollback() // err is non-nil; don't change it
			if execErr != nil {
				s.log.Error(execErr, "error exec rollback")
			}
		default:
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	return txFunc(s.WithTx(tx))
}
