package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	dbDriverName = "postgres"
	dbSource     = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	testQueries  *Queries
	testDB       *sql.DB
)

// TestMain just does extra setup and tear down for the app
func TestMain(m *testing.M) {
	var err error
	if testDB, err = sql.Open(dbDriverName, dbSource); err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
