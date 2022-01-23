package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/harrychopra/go-api/api"
	db "github.com/harrychopra/go-api/db/models"
	"github.com/harrychopra/go-api/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("failaed to create new server object: ", err)
	}

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
