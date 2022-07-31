package storage

import (
	db "github.com/ostamand/aqualog/db/sqlc"
)

type Storage interface {
	db.Querier
}
