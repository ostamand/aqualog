// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: value.sql

package db

import (
	"context"
)

const createValue = `-- name: CreateValue :one
INSERT INTO values (
  user_id,
  value_type_id,
  value
) VALUES (
  $1,
  $2,
  $3
) 
RETURNING id, user_id, value_type_id, value, created_at
`

type CreateValueParams struct {
	UserID      int64   `json:"user_id"`
	ValueTypeID int32   `json:"value_type_id"`
	Value       float64 `json:"value"`
}

func (q *Queries) CreateValue(ctx context.Context, arg CreateValueParams) (Value, error) {
	row := q.db.QueryRowContext(ctx, createValue, arg.UserID, arg.ValueTypeID, arg.Value)
	var i Value
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ValueTypeID,
		&i.Value,
		&i.CreatedAt,
	)
	return i, err
}

const listValuesPerType = `-- name: ListValuesPerType :many
SELECT id, user_id, value_type_id, value, created_at from values
WHERE value_type_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListValuesPerTypeParams struct {
	ValueTypeID int32 `json:"value_type_id"`
	Limit       int32 `json:"limit"`
	Offset      int32 `json:"offset"`
}

func (q *Queries) ListValuesPerType(ctx context.Context, arg ListValuesPerTypeParams) ([]Value, error) {
	rows, err := q.db.QueryContext(ctx, listValuesPerType, arg.ValueTypeID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Value
	for rows.Next() {
		var i Value
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ValueTypeID,
			&i.Value,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
