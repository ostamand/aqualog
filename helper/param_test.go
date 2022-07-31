package helper

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/util"
	"github.com/stretchr/testify/assert"
)

func saveRandomeUser(t *testing.T) db.User {
	user, err := SaveUser(context.Background(), store, SaveUserParams{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 6),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	return user
}

func TestCanSaveNewParam(t *testing.T) {
	user := saveRandomeUser(t)

	v := gofakeit.Float64Range(0, 10)

	value, err := SaveParam(context.Background(), store, SaveParamArgs{
		UserID:    user.ID,
		ParamName: util.GenerateRandomKey(6),
		Value:     v,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, value)
	assert.Equal(t, v, value.Value)
	assert.Equal(t, user.ID, value.UserID)
}

func TestSaveNewParamTypeExists(t *testing.T) {
	user := saveRandomeUser(t)

	paramType, err := store.CreateValueType(context.Background(), db.CreateValueTypeParams{
		Name:   util.GenerateRandomKey(6),
		UserID: user.ID,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, paramType)
	assert.Equal(t, user.ID, paramType.UserID)

	param, err := SaveParam(context.Background(), store, SaveParamArgs{
		UserID:    user.ID,
		ParamName: paramType.Name,
		Value:     gofakeit.Float64(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, param)
	assert.Equal(t, paramType.ID, param.ParamTypeID)
}
