package storage

import (
	"context"
	"database/sql"

	db "github.com/ostamand/aqualog/db/sqlc"
)

type Storage interface {
	db.Querier
	GetValueFilled(ctx context.Context, id int64) (GetValueFilledResponse, error)
	SafeCreateValue(ctx context.Context, args SafeCreateValueParams) (db.Value, error)
}

type SQLStorage struct {
	*db.Queries
	db *sql.DB
}

func NewStorage(conn *sql.DB) Storage {
	return &SQLStorage{
		db:      conn,
		Queries: db.New(conn),
	}
}

func (s SQLStorage) executeTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		// rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}

type GetValueFilledResponse struct {
	Value db.Value     `json:"value"`
	Type  db.ValueType `json:"type"`
	User  db.User      `json:"user"`
}

func (s SQLStorage) GetValueFilled(ctx context.Context, id int64) (GetValueFilledResponse, error) {
	res := GetValueFilledResponse{}

	// get value
	value, err := s.Queries.GetValue(ctx, id)
	if err != nil {
		return res, err
	}

	// get value type
	var valueType db.ValueType
	valueType, err = s.Queries.GetValueType(ctx, int64(value.ValueTypeID))
	if err != nil {
		return res, err
	}

	// get user
	var user db.User
	user, err = s.Queries.GetUser(ctx, value.UserID)
	if err != nil {
		return res, err
	}

	res.Value = value
	res.Type = valueType
	res.User = user

	return res, err
}

type SafeCreateValueParams struct {
	Username string  `json:"username"`
	Value    float64 `json:"value"`
	Type     string  `json:"type"`
}

func (s SQLStorage) SafeCreateValue(ctx context.Context, args SafeCreateValueParams) (db.Value, error) {
	var v db.Value

	err := s.executeTx(ctx, func(q *db.Queries) error {
		var err error

		// check if user first
		var user db.User
		user, err = q.GetByUsername(ctx, args.Username)
		if err != nil {
			// create new username
			user, err = q.CreateUser(ctx, args.Username)
			if err != nil {
				// TODO: check what caused error. most probably username already exists
				return err
			}
		}

		// check if value type exists
		var valueType db.ValueType
		valueType, err = q.GetValueTypeByName(ctx, db.GetValueTypeByNameParams{
			UserID: user.ID,
			Name:   args.Type,
		})
		if err != nil {
			// create new value type for user
			valueType, err = q.CreateValueType(ctx, db.CreateValueTypeParams{
				Name:   args.Type,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}
		}

		// create new measurement
		v, err = q.CreateValue(ctx, db.CreateValueParams{
			UserID:      user.ID,
			ValueTypeID: int32(valueType.ID),
			Value:       args.Value,
		})

		return err
	})

	return v, err
}
