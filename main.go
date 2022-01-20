package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/harrychopra/go-api/api"
	db "github.com/harrychopra/go-api/db/models"
)

var build = "develop"

const (
	dbDriverName  = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriverName, dbSource)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(serverAddress); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
