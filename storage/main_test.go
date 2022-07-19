package storage

import (
	"database/sql"
	"os"
	"testing"

	db "github.com/ostamand/aqualog/db/sqlc"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	testQueries, testDb = db.SetupDB()
	os.Exit(m.Run())
}
