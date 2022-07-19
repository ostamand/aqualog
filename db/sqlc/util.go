package db

import (
	"database/sql"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/aqualog?sslmode=disable"
)

func SetupDB() (*Queries, *sql.DB) {
	testDb, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	testQueries := New(testDb)
	return testQueries, testDb
}
