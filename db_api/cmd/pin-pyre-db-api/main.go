package main

import (
	"log"

	"github.com/RockOnRoad/pin_pyre/db_api/api"
	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

func main() {
	sqlxDB, err := db.NewSQLite("data/tyres.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlxDB.Close()

	if err := db.Migrate(sqlxDB, "./migrations"); err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(sqlxDB)
	server := api.NewServer(store)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
