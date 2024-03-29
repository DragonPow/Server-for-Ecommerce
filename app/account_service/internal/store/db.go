// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package store

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.getUserByUserNameAndPasswordStmt, err = db.PrepareContext(ctx, getUserByUserNameAndPassword); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUserNameAndPassword: %w", err)
	}
	if q.getAccountByIdsStmt, err = db.PrepareContext(ctx, getAccountByIds); err != nil {
		return nil, fmt.Errorf("error preparing query getAccountByIds: %w", err)
	}
	if q.getCustomerByIdsStmt, err = db.PrepareContext(ctx, getCustomerByIds); err != nil {
		return nil, fmt.Errorf("error preparing query getCustomerByIds: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.getUserByUserNameAndPasswordStmt != nil {
		if cerr := q.getUserByUserNameAndPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByUserNameAndPasswordStmt: %w", cerr)
		}
	}
	if q.getAccountByIdsStmt != nil {
		if cerr := q.getAccountByIdsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountByIdsStmt: %w", cerr)
		}
	}
	if q.getCustomerByIdsStmt != nil {
		if cerr := q.getCustomerByIdsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCustomerByIdsStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                               DBTX
	tx                               *sql.Tx
	getUserByUserNameAndPasswordStmt *sql.Stmt
	getAccountByIdsStmt              *sql.Stmt
	getCustomerByIdsStmt             *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                               tx,
		tx:                               tx,
		getUserByUserNameAndPasswordStmt: q.getUserByUserNameAndPasswordStmt,
		getAccountByIdsStmt:              q.getAccountByIdsStmt,
		getCustomerByIdsStmt:             q.getCustomerByIdsStmt,
	}
}
