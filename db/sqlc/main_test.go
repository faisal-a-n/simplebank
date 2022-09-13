package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/faisal-a-n/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("Coudln't load config %v", err)
	}

	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
