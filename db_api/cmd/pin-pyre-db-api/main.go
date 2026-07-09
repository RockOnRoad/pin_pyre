package main

import (
	"log"

	"github.com/RockOnRoad/pin_pyre/db_api/api"
	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

func main() {
	sqlxDB, err := db.NewSQLite("db/tyres.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlxDB.Close()

	store := db.NewStore(sqlxDB)
	server := api.NewServer(store)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
