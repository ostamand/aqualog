package db

import (
	"context"
	"database/sql"
)

type Storage struct {
	*Queries
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db:      db,
		Queries: New(db),
	}
}

func (s Storage) executeTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
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

type AddMeasurementParams struct {
	Username string  `json:"username"`
	Value    float64 `json:"value"`
	Type     string  `json:"type"`
}

func (s Storage) AddMeasurement(ctx context.Context, args AddMeasurementParams) (Value, error) {
	var v Value

	err := s.executeTx(ctx, func(q *Queries) error {
		var err error

		// check if user first
		var user User
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
		var valueType ValueType
		valueType, err = q.GetValueTypeByName(ctx, GetValueTypeByNameParams{
			UserID: user.ID,
			Name:   args.Type,
		})
		if err != nil {
			// create new value type for user
			valueType, err = q.CreateValueType(ctx, CreateValueTypeParams{
				Name:   args.Type,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}
		}

		// create new measurement
		v, err = q.CreateValue(ctx, CreateValueParams{
			UserID:      user.ID,
			ValueTypeID: int32(valueType.ID),
			Value:       args.Value,
		})

		return err
	})

	return v, err
}
