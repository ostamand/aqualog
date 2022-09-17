package helper

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/storage"
	"github.com/ostamand/aqualog/util"
	"github.com/stretchr/testify/assert"
)

func TestSummary(t *testing.T) {
	ctx := context.Background()
	user := SaveRandomUser(t)
	now := time.Now()

	data := []SaveParamArgs{
		{
			ParamName: "nitrate",
			Value:     12,
			Timestamp: now,
		},
		{
			ParamName: "nitrate",
			Value:     24,
			Timestamp: now.Add(time.Hour * 5),
		},
		{
			ParamName: "phosphate",
			Value:     0.024,
			Timestamp: now.Add(-time.Hour * 5),
		},
	}

	for _, item := range data {
		value, err := SaveParam(ctx, store, SaveParamArgs{
			UserID:    user.ID,
			ParamName: item.ParamName,
			Value:     item.Value,
			Timestamp: item.Timestamp,
		})
		assert.NoError(t, err)
		assert.Equal(t, item.Value, value.Value)
		// assert.Equal(t, item.Timestamp, value.Timestamp) TODO: time.Local VS time.UTC
		assert.Equal(t, user.ID, value.UserID)
	}

	items, err := store.ListSummary(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, len(data)-1, len(items))

	findWithName := func(name string) *storage.ListSummaryRow {
		// find & check nitrate
		var line *storage.ListSummaryRow
		for _, item := range items {
			if item.Name == name {
				line = &item
				break
			}
		}
		return line
	}

	// check nitrate
	line := findWithName(data[0].ParamName)
	assert.NotNil(t, line)
	assert.Equal(t, data[1].Value, line.Value)
	assert.Equal(t, data[0].Value, *line.PreviousValue)

	// check phosphate
	line = findWithName(data[2].ParamName)
	assert.NotNil(t, line)
	assert.Equal(t, data[2].Value, line.Value)
	assert.Nil(t, line.PreviousValue)
	assert.Nil(t, line.PreviousTimestamp)
}

func TestCanSaveNewParam(t *testing.T) {
	user := SaveRandomUser(t)

	v := gofakeit.Float64Range(0, 10)

	value, err := SaveParam(context.Background(), store, SaveParamArgs{
		UserID:    user.ID,
		ParamName: util.GenerateRandomString(6),
		Value:     v,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, value)
	assert.Equal(t, v, value.Value)
	assert.Equal(t, user.ID, value.UserID)
}

func TestSaveNewParamTypeExists(t *testing.T) {
	user := SaveRandomUser(t)

	paramType, err := store.CreateParamType(context.Background(), db.CreateParamTypeParams{
		Name:   util.GenerateRandomString(6),
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

func TestCanCreateParams(t *testing.T) {
	now := time.Now()

	// create random user
	user, err := testQueries.CreateUser(context.Background(), db.CreateUserParams{
		Username:       gofakeit.Username(),
		Email:          gofakeit.Email(),
		HashedPassword: util.GenerateRandomString(6),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	//create random param type
	paramType, err := testQueries.CreateParamType(context.Background(), db.CreateParamTypeParams{
		Name:   util.GenerateRandomString(6),
		UserID: user.ID,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, paramType)
	assert.Equal(t, user.ID, paramType.UserID)

	// define test cases
	testCases := []struct {
		name   string
		args   SaveParamArgs
		expect time.Time
	}{
		{
			name: "No timestamp",
			args: SaveParamArgs{
				UserID:    user.ID,
				ParamName: paramType.Name,
			},
			expect: now,
		},
		{
			name: "Timstamp provided",
			args: SaveParamArgs{
				UserID:    user.ID,
				ParamName: paramType.Name,
				Timestamp: now.Add(-5 * time.Hour),
			},
			expect: now.Add(-5 * time.Hour),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			param, err := SaveParam(context.Background(), store, tc.args)

			assert.NoError(t, err)
			assert.NotEmpty(t, param)
			assert.Equal(t, tc.args.UserID, param.UserID)
			assert.Equal(t, paramType.ID, param.ParamTypeID)

			assert.WithinDuration(t, tc.expect, param.Timestamp, time.Second)
		})
	}
}

func TestGetParams(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name   string
		times  []time.Time
		from   time.Time
		to     time.Time
		expect int
	}{
		{
			name: "All params with type",
			times: []time.Time{
				now.Add(-time.Minute * 2),
				now.Add(time.Hour),
				now.Add(time.Hour * 24),
			},
			expect: 3,
		},
		{
			name: "Params with from & to",
			times: []time.Time{
				now.Add(-time.Minute * 5),
				now.Add(time.Minute * 1),
				now.Add(time.Minute * 2),
				now.Add(time.Minute * 10),
				now.Add(time.Minute * 20),
			},
			from:   now,
			to:     now.Add(time.Minute * 10),
			expect: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create random user
			user := SaveRandomUser(t)

			// create random param type
			paramName := util.GenerateRandomString(6)

			// create params
			for _, tick := range tc.times {
				_, err := SaveParam(context.Background(), store, SaveParamArgs{
					UserID:    user.ID,
					ParamName: paramName,
					Value:     gofakeit.Float64(),
					Timestamp: tick,
				})
				assert.NoError(t, err)
			}

			// query params using helper
			args := db.ListParamsByTypeParams{
				UserID:        user.ID,
				ParamTypeName: paramName,
				From:          tc.from,
				To:            tc.to,
			}
			args.FillDefaults()

			params, err := store.ListParamsByType(context.Background(), args)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, len(params))
		})
	}
}
