package storage

import (
	"context"

	db "github.com/ostamand/aqualog/db/sqlc"
)

type Storage interface {
	db.Querier
	ListSummary(ctx context.Context, userID int64) ([]ListSummaryRow, error)
}
