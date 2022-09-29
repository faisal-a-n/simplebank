package token

import (
	"testing"
	"time"

	"github.com/faisal-a-n/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateString(32))
	require.NoError(t, err)

	id := util.GenerateRandomInt(1000, 1)
	duration := time.Minute
	issued_at := time.Now()
	expires_at := time.Now().Add(time.Minute)

	token, _, err := maker.CreateToken(id, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issued_at, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expires_at, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateString(32))
	require.NoError(t, err)

	id := util.GenerateRandomInt(1000, 1)

	token, _, err := maker.CreateToken(id, -1)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ERR_TOKEN_EXPIRED.Error())
	require.Nil(t, payload)
}
