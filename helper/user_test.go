package helper

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveUser(t *testing.T) {
	args := SaveUserParams{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 8),
	}
	user, err := SaveUser(context.Background(), s, args)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
}

func TestPasswordTooShort(t *testing.T) {
	args := SaveUserParams{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 5),
	}
	user, err := SaveUser(context.Background(), s, args)
	assert.ErrorIs(t, err, ErrPasword{Info: PasswordLength})
	assert.Empty(t, user)
}
