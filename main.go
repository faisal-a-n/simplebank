package main

import (
	"database/sql"
	"log"

	"github.com/faisal-a-n/simplebank/api"
	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/faisal-a-n/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Coudln't load config %v", err)
	}

	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatalf("Coudln't connect to db %v", err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.PORT); err != nil {
		log.Fatalf("Coudln't start server %v", err.Error())
	}
}
