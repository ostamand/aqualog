package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	userName := gofakeit.Username()
	user, err := testQueries.CreateUser(context.Background(), userName)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, userName)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)

	assert.NoError(t, err)
	assert.Equal(t, user1.Username, user2.Username)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	assert.NoError(t, err)
	_, err = testQueries.GetUser(context.Background(), user.ID)
	assert.Error(t, err)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}
	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
