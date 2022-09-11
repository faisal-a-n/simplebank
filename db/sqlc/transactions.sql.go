// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: transactions.sql

package db

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT into transactions (
 "from_account_id", "to_account_id", "from_entry_id", "to_entry_id", "amount", "created_at"
)
values
($1, $2, $3, $4, $5, $6) RETURNING id, from_account_id, to_account_id, from_entry_id, to_entry_id, amount, created_at
`

type CreateTransactionParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	FromEntryID   int64 `json:"from_entry_id"`
	ToEntryID     int64 `json:"to_entry_id"`
	Amount        int64 `json:"amount"`
	CreatedAt     int64 `json:"created_at"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.FromEntryID,
		arg.ToEntryID,
		arg.Amount,
		arg.CreatedAt,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.FromEntryID,
		&i.ToEntryID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, from_account_id, to_account_id, from_entry_id, to_entry_id, amount, created_at from transactions where id = $1 limit 1
`

func (q *Queries) GetTransaction(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.FromEntryID,
		&i.ToEntryID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listTransactions = `-- name: ListTransactions :many
SELECT id, from_account_id, to_account_id, from_entry_id, to_entry_id, amount, created_at from transactions order by id limit $1 offset $2
`

type ListTransactionsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, listTransactions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.FromEntryID,
			&i.ToEntryID,
			&i.Amount,
			&i.CreatedAt,
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
