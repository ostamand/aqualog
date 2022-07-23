package token

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
)

const keySize = chacha20poly1305.KeySize

func generateRandomKey(size int) string {
	chars := "abcdefghijklmnopqrtuvwxyz"
	n := len(chars)
	var sb strings.Builder
	for i := 0; i < size; i++ {
		sb.WriteByte(chars[rand.Intn(n)])
	}
	return sb.String()
}

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(generateRandomKey(keySize))
	require.NoError(t, err)

	username := gofakeit.Username()

	token, _, err := maker.CreateToken(username, time.Minute*10)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, username, payload.Username)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(generateRandomKey(keySize))
	require.NoError(t, err)

	token, _, err := maker.CreateToken(gofakeit.Username(), -time.Minute*10)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}

func TestInvalidKeySize(t *testing.T) {
	maker, err := NewPasetoMaker(generateRandomKey(12))
	require.EqualError(t, err, ErrInvalidKeySize.Error())
	require.Empty(t, maker)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(generateRandomKey(keySize))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	token, payload, err := maker.CreateToken(gofakeit.Username(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	otherMaker, err := NewPasetoMaker(generateRandomKey(keySize))
	require.NoError(t, err)
	require.NotEmpty(t, otherMaker)

	payload, err = otherMaker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, payload)
}
