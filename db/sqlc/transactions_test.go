package db

import (
	"context"
	"testing"
	"time"

	"github.com/faisal-a-n/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestTransaction(t *testing.T) Transaction {
	amount := util.GenerateAmount()
	entry1, account1 := createTestEntry(t, -amount)

	entry2, account2 := createTestEntry(t, amount)

	args := CreateTransactionParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		FromEntryID:   entry1.ID,
		ToEntryID:     entry2.ID,
		Amount:        amount,
		CreatedAt:     time.Now().Unix(),
	}
	tx, err := testQueries.CreateTransaction(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, tx.FromAccountID, account1.ID)
	require.Equal(t, tx.ToAccountID, account2.ID)
	require.Equal(t, tx.FromEntryID, entry1.ID)
	require.Equal(t, tx.ToEntryID, entry2.ID)
	require.Equal(t, tx.Amount, entry2.Amount)
	require.Equal(t, -tx.Amount, entry1.Amount)
	return tx
}

func TestCreateTransaction(t *testing.T) {
	createTestTransaction(t)
}

func TestGetTransaction(t *testing.T) {
	tx := createTestTransaction(t)
	checkTx, err := testQueries.GetTransaction(context.Background(), tx.ID)
	require.NoError(t, err)
	require.Equal(t, tx.ID, checkTx.ID)
	require.Equal(t, tx.Amount, checkTx.Amount)
	require.Equal(t, tx.FromAccountID, checkTx.FromAccountID)
	require.Equal(t, tx.ToAccountID, checkTx.ToAccountID)
}

func TestListTransactions(t *testing.T) {
	for i := 0; i < 5; i++ {
		createTestTransaction(t)
	}

	args := ListTransactionsParams{
		Limit:  100,
		Offset: 0,
	}

	list, err := testQueries.ListTransactions(context.Background(), args)

	require.NoError(t, err)

	for _, entry := range list[len(list)-5:] {
		require.NotEmpty(t, entry)
	}
}
