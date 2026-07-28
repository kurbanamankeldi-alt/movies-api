package main

import (
	"log"

	"github.com/kurbanamankeldi-alt/movies-api/db"

	"github.com/kurbanamankeldi-alt/movies-api/server

)

func main() {
	//initialize database
	database, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//seed with sample data
	err = db.SeedTables(database)
	if err != nil {
		panic(err)
	}
	//call for server
	srv := server.Server(database)
	log.Fatal(srv.ListenAndServe())
}
