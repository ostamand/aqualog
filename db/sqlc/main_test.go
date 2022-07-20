package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	testQueries, testDb = SetupDB("../..")
	os.Exit(m.Run())
}
