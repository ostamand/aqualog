package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ostamand/aqualog/api"
	"github.com/ostamand/aqualog/storage"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/aqualog?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	s := storage.NewStorage(conn)
	server := api.NewServer(s)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
