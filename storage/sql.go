package storage

import (
	"context"
	"database/sql"
	"time"

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

type ListSummaryRow struct {
	Name          string     `json:"name"`
	Value         float64    `json:"value"`
	Timestamp     time.Time  `json:"timestamp"`
	LastValue     *float64   `json:"last_value"`
	LastTimestamp *time.Time `json:"last_timestamp"`
}

const listSummary = `WITH a AS (
	SELECT 
	DISTINCT ON (params.param_type_id) param_type_id,
	params.id,
	t."name",
	t.target,
	t."min",
	t."max",
	params."value",
	params.timestamp,
	params.created_at
	FROM params
	INNER JOIN param_types as t ON params.param_type_id = t.id
	WHERE params.user_id = $1
	ORDER BY param_type_id, "timestamp" DESC
)
SELECT
a.name,
a."value",
a."timestamp",
b.last_value,
b.last_timestamp
FROM a 
LEFT JOIN (
	SELECT
	DISTINCT ON (params.param_type_id) params.param_type_id,
	params."value" as last_value,
	params.timestamp as last_timestamp
	FROM params
	RIGHT JOIN a ON params.param_type_id = a.param_type_id
	WHERE params.user_id = $1 AND
		params.id NOT IN (SELECT id FROM a)
	ORDER BY param_type_id, params."timestamp" DESC
) as b
ON a.param_type_id = b.param_type_id;`

func (s SQLStorage) ListSummary(ctx context.Context, userID int64) ([]ListSummaryRow, error) {
	rows, err := s.db.QueryContext(ctx, listSummary, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSummaryRow{}
	for rows.Next() {
		var row ListSummaryRow
		if err := rows.Scan(&row.Name, &row.Value, &row.Timestamp, &row.LastValue, &row.LastTimestamp); err != nil {
			return nil, err
		}
		items = append(items, row)
	}
	return items, nil
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
