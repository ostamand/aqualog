package helper

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/ostamand/aqualog/db/sqlc"
	"github.com/ostamand/aqualog/storage"
)

var testQueries *db.Queries
var testDb *sql.DB
var store storage.Storage

func TestMain(m *testing.M) {
	testQueries, testDb = db.SetupDB("../")
	store = storage.NewSQLStorage(testDb)
	os.Exit(m.Run())
}
