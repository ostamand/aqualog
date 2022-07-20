package db

import (
	"database/sql"
	"log"

	"github.com/ostamand/aqualog/util"
)

func SetupDB(configPath string) (*Queries, *sql.DB) {
	config, err := util.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Cannot load configs", err)
	}
	testDb, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	testQueries := New(testDb)
	return testQueries, testDb
}
