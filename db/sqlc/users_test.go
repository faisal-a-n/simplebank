package db

import (
	"context"
	"testing"
	"time"

	"github.com/faisal-a-n/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) User {
	hash, err := util.HashPassword(util.GenerateString(12))
	require.NoError(t, err)

	arg := CreateUserParams{
		Name:              util.GenerateString(8),
		Email:             util.RandomEmail(),
		Password:          hash,
		PasswordChangedAt: time.Now().Unix(),
		CreatedAt:         time.Now().Unix(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {
	user := createTestUser(t)
	fetchedUser, err := testQueries.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.Equal(t, user.ID, fetchedUser.ID)
	require.Equal(t, user.Email, fetchedUser.Email)
	require.Equal(t, user.Password, fetchedUser.Password)
	require.Equal(t, user.Name, fetchedUser.Name)
}
