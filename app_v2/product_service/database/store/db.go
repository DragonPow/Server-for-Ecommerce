// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

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
	if q.getCategoriesStmt, err = db.PrepareContext(ctx, getCategories); err != nil {
		return nil, fmt.Errorf("error preparing query GetCategories: %w", err)
	}
	if q.getProductAndRelationsStmt, err = db.PrepareContext(ctx, getProductAndRelations); err != nil {
		return nil, fmt.Errorf("error preparing query GetProductAndRelations: %w", err)
	}
	if q.getProductDetailsStmt, err = db.PrepareContext(ctx, getProductDetails); err != nil {
		return nil, fmt.Errorf("error preparing query GetProductDetails: %w", err)
	}
	if q.getProductTemplatesStmt, err = db.PrepareContext(ctx, getProductTemplates); err != nil {
		return nil, fmt.Errorf("error preparing query GetProductTemplates: %w", err)
	}
	if q.getProductsStmt, err = db.PrepareContext(ctx, getProducts); err != nil {
		return nil, fmt.Errorf("error preparing query GetProducts: %w", err)
	}
	if q.getProductsByKeywordStmt, err = db.PrepareContext(ctx, getProductsByKeyword); err != nil {
		return nil, fmt.Errorf("error preparing query GetProductsByKeyword: %w", err)
	}
	if q.getSellersStmt, err = db.PrepareContext(ctx, getSellers); err != nil {
		return nil, fmt.Errorf("error preparing query GetSellers: %w", err)
	}
	if q.getUomsStmt, err = db.PrepareContext(ctx, getUoms); err != nil {
		return nil, fmt.Errorf("error preparing query GetUoms: %w", err)
	}
	if q.getUsersStmt, err = db.PrepareContext(ctx, getUsers); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsers: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.getCategoriesStmt != nil {
		if cerr := q.getCategoriesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCategoriesStmt: %w", cerr)
		}
	}
	if q.getProductAndRelationsStmt != nil {
		if cerr := q.getProductAndRelationsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProductAndRelationsStmt: %w", cerr)
		}
	}
	if q.getProductDetailsStmt != nil {
		if cerr := q.getProductDetailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProductDetailsStmt: %w", cerr)
		}
	}
	if q.getProductTemplatesStmt != nil {
		if cerr := q.getProductTemplatesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProductTemplatesStmt: %w", cerr)
		}
	}
	if q.getProductsStmt != nil {
		if cerr := q.getProductsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProductsStmt: %w", cerr)
		}
	}
	if q.getProductsByKeywordStmt != nil {
		if cerr := q.getProductsByKeywordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProductsByKeywordStmt: %w", cerr)
		}
	}
	if q.getSellersStmt != nil {
		if cerr := q.getSellersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSellersStmt: %w", cerr)
		}
	}
	if q.getUomsStmt != nil {
		if cerr := q.getUomsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUomsStmt: %w", cerr)
		}
	}
	if q.getUsersStmt != nil {
		if cerr := q.getUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsersStmt: %w", cerr)
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
	db                         DBTX
	tx                         *sql.Tx
	getCategoriesStmt          *sql.Stmt
	getProductAndRelationsStmt *sql.Stmt
	getProductDetailsStmt      *sql.Stmt
	getProductTemplatesStmt    *sql.Stmt
	getProductsStmt            *sql.Stmt
	getProductsByKeywordStmt   *sql.Stmt
	getSellersStmt             *sql.Stmt
	getUomsStmt                *sql.Stmt
	getUsersStmt               *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                         tx,
		tx:                         tx,
		getCategoriesStmt:          q.getCategoriesStmt,
		getProductAndRelationsStmt: q.getProductAndRelationsStmt,
		getProductDetailsStmt:      q.getProductDetailsStmt,
		getProductTemplatesStmt:    q.getProductTemplatesStmt,
		getProductsStmt:            q.getProductsStmt,
		getProductsByKeywordStmt:   q.getProductsByKeywordStmt,
		getSellersStmt:             q.getSellersStmt,
		getUomsStmt:                q.getUomsStmt,
		getUsersStmt:               q.getUsersStmt,
	}
}