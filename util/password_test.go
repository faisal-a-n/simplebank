package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := GenerateString(12)
	hash, err := HashPassword(password)
	require.NoError(t, err)

	hash2, err := HashPassword(password)
	require.NoError(t, err)

	require.NotEqual(t, hash, hash2)

	err = CheckPassword(password, hash)
	require.NoError(t, err)

	password = GenerateString(10)
	err = CheckPassword(password, hash)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
