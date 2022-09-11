package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createTestAccount(t, -1)
	account2 := createTestAccount(t, -1)

	fmt.Println("before >> ", account1.Balance, account2.Balance)
	//Run transfers with go routines

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//Check transaction
		transaction := result.Transaction
		require.NotEmpty(t, transaction)
		require.Equal(t, transaction.FromAccountID, account1.ID)
		require.Equal(t, transaction.ToAccountID, account2.ID)
		require.Equal(t, transaction.Amount, amount)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.CreatedAt)

		_, err = store.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		//Check account entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount.Balance)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount.Balance)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check account balance
		fmt.Println("tx >> ", fromAccount.Balance, toAccount.Balance)

		moneySent := account1.Balance - fromAccount.Balance
		moneyReceived := toAccount.Balance - account2.Balance

		require.True(t, moneySent > 0)
		require.Equal(t, moneySent, moneyReceived)
		require.True(t, moneySent%amount == 0)

		// Since the transactions are being done on a single account the balance keeps decreasing n*amount times
		k := int(moneySent / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//Check the final balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Println("after >> ", updatedAccount1.Balance, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	var err error
	account1 := createTestAccount(t, -1)
	account2 := createTestAccount(t, -1)

	fmt.Println("before >> ", account1.Balance, account2.Balance)
	//Run transfers with go routines

	n := 5
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			//Queries are ran sequentially to prevent deadlock. queries with same id are ran together so they don't have to wait
			if account1.ID < account2.ID {
				account1, account2, err = addMoney(context.Background(), store.Queries, account1.ID, account2.ID, -amount, amount)
			} else {
				account1, account2, err = addMoney(context.Background(), store.Queries, account2.ID, account1.ID, amount, -amount)
			}
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	//Check the final balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	fmt.Println("after >> ", updatedAccount1.Balance, updatedAccount2.Balance)
}

func addMoney(ctx context.Context, q *Queries, account1Id, account2Id, amount1, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateBalance(ctx, UpdateBalanceParams{
		Amount: amount1,
		ID:     account1Id,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateBalance(ctx, UpdateBalanceParams{
		Amount: amount2,
		ID:     account2Id,
	})
	if err != nil {
		return
	}
	return
}
