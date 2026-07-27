package main

import (
	"fmt"
    "log"
    "net/http"	
	"github.com/kurbanamankeldi-alt/movies-api/db"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
	"github.com/kurbanamankeldi-alt/movies-api/service"
	"github.com/kurbanamankeldi-alt/movies-api/handler"
	"github.com/kurbanamankeldi-alt/movies-api/errors"
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
	mux.Handle("GET /api/movie", errors.HttpErrorHandler(handler.Get))
	mux.Handle("GET /api/movie/{id}", errors.HttpErrorHandler(handler.GetById))
	mux.Handle("POST /api/movie", errors.HttpErrorHandler(handler.Create))
	fmt.Println("Server is running in the below link:")
	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))	
}