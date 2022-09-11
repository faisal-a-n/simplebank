package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store provides functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// Create new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//Executes db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("Tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

var txKey = struct{}{}

//Input for transfer tx
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	ToEntryID     int64 `json:"to_entry_id"`
	FromEntryID   int64 `json:"from_entry_id"`
	Amount        int64 `json:"amount"`
	CreatedAt     int64 `json:"created_at"`
}

//Result of transfer tx
type TransferTxResult struct {
	Transaction Transaction `json:"transfer"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
	FromEntry   Entry       `json:"from_entry"`
	ToEntry     Entry       `json:"to_entry"`
}

//This function is transfers money between accounts and adds entries to entries table, all done in single transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
			CreatedAt: time.Now().Unix(),
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
			CreatedAt: time.Now().Unix(),
		})

		if err != nil {
			return err
		}

		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			FromEntryID:   result.FromEntry.ID,
			ToEntryID:     result.ToEntry.ID,
			Amount:        arg.Amount,
			CreatedAt:     time.Now().Unix(),
		})

		if err != nil {
			return err
		}

		account1, err := q.UpdateBalance(ctx, UpdateBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromAccount = account1

		account2, err := q.UpdateBalance(ctx, UpdateBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToAccount = account2

		return nil
	})

	return result, err
}
