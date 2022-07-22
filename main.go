package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ostamand/aqualog/api"
	"github.com/ostamand/aqualog/storage"
	"github.com/ostamand/aqualog/util"
)

func main() {
	// load configs
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configs", err)
	}

	// connect to db
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}

	s := storage.NewSQLStorage(conn)

	// start server
	server := api.NewServer(s)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
