package main

import (
	"fmt"
    "log"
    "net/http"	
	"github.com/kurbanamankeldi-alt/movies-api/db"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
	"github.com/kurbanamankeldi-alt/movies-api/service"
	"github.com/kurbanamankeldi-alt/movies-api/handler"
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

	port := "8081"
	repo := repository.NewSQLiteMovieRepository(database)
	service := service.NewMovieService(repo)
	handler := handler.NewMovieHandler(service)	
	mux := http.NewServeMux()
	mux.HandleFunc("/api/movie/", handler.GetById)	
	mux.HandleFunc("/api/movie", handler.CreateMovie)
	fmt.Println("Server is running in the below link:")
	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))	
}