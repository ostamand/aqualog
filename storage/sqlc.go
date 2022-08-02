package storage

import (
	"context"
	"database/sql"

	db "github.com/ostamand/aqualog/db/sqlc"
)

type SQLStorage struct {
	*db.Queries
	db *sql.DB
}

func NewSQLStorage(conn *sql.DB) Storage {
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
