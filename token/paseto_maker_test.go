package token

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
)

const keySize = chacha20poly1305.KeySize

func generateRandomUser() db.User {
	return db.User{
		ID:       gofakeit.Int64(),
		Username: gofakeit.Username(),
	}
}

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateRandomKey(keySize))
	require.NoError(t, err)

	user := generateRandomUser()

	token, _, err := maker.CreateToken(CreateTokenArgs{
		Username: user.Username,
		UserID:   user.ID,
		Duration: time.Minute * 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, user.Username, payload.Username)
	require.Equal(t, user.ID, payload.UserID)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateRandomKey(keySize))
	require.NoError(t, err)

	user := generateRandomUser()

	token, _, err := maker.CreateToken(CreateTokenArgs{
		Username: user.Username,
		UserID:   user.ID,
		Duration: -time.Minute * 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}

func TestInvalidKeySize(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateRandomKey(12))
	require.ErrorIs(t, err, ErrInvalidKeySize)
	require.Empty(t, maker)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.GenerateRandomKey(keySize))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	user := generateRandomUser()

	token, payload, err := maker.CreateToken(CreateTokenArgs{
		Username: user.Username,
		UserID:   user.ID,
		Duration: time.Minute,
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	otherMaker, err := NewPasetoMaker(util.GenerateRandomKey(keySize))
	require.NoError(t, err)
	require.NotEmpty(t, otherMaker)

	payload, err = otherMaker.VerifyToken(token)
	require.ErrorIs(t, err, ErrInvalidToken)
	require.Empty(t, payload)
}
