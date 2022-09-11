package db

import (
	"context"
	"testing"
	"time"

	"github.com/faisal-a-n/util"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T, amount int64) (Entry, Account) {
	account := createTestAccount(t, amount)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
		CreatedAt: time.Now().Unix(),
	}
	require.LessOrEqual(t, args.Amount, account.Balance)
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.Equal(t, args.CreatedAt, entry.CreatedAt)
	return entry, account
}

func TestCreateEntry(t *testing.T) {
	createTestEntry(t, util.GenerateAmount())
}

func TestGetEntry(t *testing.T) {
	entry, _ := createTestEntry(t, util.GenerateAmount())
	checkEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.Equal(t, entry.ID, checkEntry.ID)
	require.Equal(t, entry.Amount, checkEntry.Amount)
	require.Equal(t, entry.AccountID, checkEntry.AccountID)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		createTestEntry(t, util.GenerateAmount())
	}

	args := ListEntriesParams{
		Limit:  100,
		Offset: 0,
	}

	list, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)

	for _, entry := range list[len(list)-5:] {
		require.NotEmpty(t, entry)
	}
}
