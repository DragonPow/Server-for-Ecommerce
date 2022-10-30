// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: account.sql

package store

import (
	"context"
)

const getUserByUserNameAndPassword = `-- name: GetUserByUserNameAndPassword :one
SELECT id, username, password, create_date, write_date
FROM account
WHERE username = $1::varchar
AND password = $2::varchar
LIMIT 1
`

type GetUserByUserNameAndPasswordParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) GetUserByUserNameAndPassword(ctx context.Context, arg GetUserByUserNameAndPasswordParams) (Account, error) {
	row := q.queryRow(ctx, q.getUserByUserNameAndPasswordStmt, getUserByUserNameAndPassword, arg.Username, arg.Password)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreateDate,
		&i.WriteDate,
	)
	return i, err
}
