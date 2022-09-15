package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/faisal-a-n/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T, balance int64) Account {
	user := createTestUser(t)
	if balance == -1 {
		balance = util.GenerateAmount()
	}
	arg := CreateAccountParams{
		Name:      user.Name,
		UserID:    user.ID,
		Balance:   balance,
		Currency:  util.GenerateCurrency(),
		CreatedAt: time.Now().Unix(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t, -1)
}

func TestGetAccount(t *testing.T) {
	account := createTestAccount(t, -1)
	account2, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account.ID, account2.ID)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createTestAccount(t, -1)
	}

	args := ListAccountsParams{
		Limit:  100,
		Offset: 0,
	}

	list, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)

	for _, account := range list[len(list)-5:] {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createTestAccount(t, -1)
	amount := util.GenerateAmount()

	args := UpdateBalanceParams{
		ID:     account.ID,
		Amount: amount,
	}

	uAccount, err := testQueries.UpdateBalance(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, args.Amount, uAccount.Balance-account.Balance)
	require.Equal(t, account.ID, uAccount.ID)
}

func TestDeleteAccount(t *testing.T) {
	account := createTestAccount(t, -1)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	checkAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, checkAccount)
}

func TestListAccountsForUser(t *testing.T) {
	account := createTestAccount(t, -1)

	args := ListAccountsForUserParams{
		UserID: account.UserID,
		Limit:  100,
		Offset: 0,
	}

	list, err := testQueries.ListAccountsForUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, list)
	require.Equal(t, account.UserID, list[0].UserID)
}
