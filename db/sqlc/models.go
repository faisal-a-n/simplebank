// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import ()

type Account struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
	CreatedAt int64  `json:"created_at"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
	CreatedAt int64 `json:"created_at"`
}

type Transaction struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	FromEntryID   int64 `json:"from_entry_id"`
	ToEntryID     int64 `json:"to_entry_id"`
	Amount        int64 `json:"amount"`
	CreatedAt     int64 `json:"created_at"`
}