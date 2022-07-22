package storage

import (
	"context"

	db "github.com/ostamand/aqualog/db/sqlc"
)

type GetValueFilledResponse struct {
	Value db.Value     `json:"value"`
	Type  db.ValueType `json:"type"`
	User  db.User      `json:"user"`
}

type SafeCreateValueParams struct {
	Username string  `json:"username"`
	Value    float64 `json:"value"`
	Type     string  `json:"type"`
}

type Storage interface {
	db.Querier
	GetValueFilled(ctx context.Context, id int64) (GetValueFilledResponse, error)
	SafeCreateValue(ctx context.Context, args SafeCreateValueParams) (db.Value, error)
}
