package storage

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestNewMeasurement(t *testing.T) {
	s := NewStorage(testDb)

	value := gofakeit.Float64()
	username := gofakeit.Username()
	name := gofakeit.LoremIpsumWord()

	m, err := s.AddMeasurement(context.Background(), AddMeasurementParams{
		Username: username,
		Value:    value,
		Type:     name,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, m)

	assert.Equal(t, m.Value, value)

	// check user was created
	var user db.User
	assert.NotEmpty(t, m.UserID)
	user, err = testQueries.GetUser(context.Background(), m.UserID)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)

	// check value type was created
	var valueType db.ValueType
	assert.NotEmpty(t, m.ValueTypeID)
	valueType, err = testQueries.GetValueTypeByName(context.Background(), db.GetValueTypeByNameParams{
		UserID: m.UserID,
		Name:   name,
	})
	assert.NoError(t, err)
	assert.Equal(t, name, valueType.Name)
}

func TestUniqueValueType(t *testing.T) {
	s := NewStorage(testDb)

	username := gofakeit.Username()
	valueTypeName := gofakeit.LoremIpsumWord()

	n := 5

	// save a couple of values for the same type

	var valueTypeID int32 // just to get ValueTypeID and UserID
	var userID int64

	for i := 0; i < n; i++ {
		v, err := s.AddMeasurement(context.Background(), AddMeasurementParams{
			Username: username,
			Value:    gofakeit.Float64(),
			Type:     valueTypeName,
		})
		assert.NoError(t, err)
		valueTypeID = v.ValueTypeID
		userID = v.UserID
	}

	// check that we have n values for value type
	values, err := testQueries.ListValuesPerType(context.Background(), db.ListValuesPerTypeParams{
		ValueTypeID: valueTypeID,
		Limit:       int32(2 * n),
		Offset:      0,
	})
	assert.NoError(t, err)
	assert.Len(t, values, n)

	for _, v := range values {
		assert.Equal(t, valueTypeID, v.ValueTypeID)
		assert.Equal(t, userID, v.UserID)
	}
}
