package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/harrychopra/go-api/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
	testStore   *Store
)

// TestMain performs setup and tear down for the tests
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	if testDB, err = sql.Open(config.DBDriver, config.DBSource); err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	testQueries = New(testDB)

	testStore = &Store{
		Queries: testQueries,
		db:      testDB,
	}
	os.Exit(m.Run())
}
