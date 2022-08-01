package helper

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/stretchr/testify/assert"
)

func SaveRandomUser(t *testing.T) db.User {
	user, err := SaveUser(context.Background(), store, SaveUserParams{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 6),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	return user
}
