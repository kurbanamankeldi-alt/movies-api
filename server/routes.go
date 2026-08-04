package server

import (
	"database/sql"
	"net/http"

	"github.com/kurbanamankeldi-alt/movies-api/errors"
	"github.com/kurbanamankeldi-alt/movies-api/handler"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
	"github.com/kurbanamankeldi-alt/movies-api/service"
)

func RegisterRoutes(mux *http.ServeMux, database *sql.DB) {
	movieRepo := repository.NewSQLiteMovieRepository(database)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handler.NewMovieHandler(movieService)

	//actors
	actorRepo := repository.NewSQLiteActorRepository(database)
	actorService := service.NewActorService(actorRepo)
	actorHandler := handler.NewActorHandler(actorService)

	mux.Handle("GET /api/movie", errors.HttpErrorHandler(movieHandler.Get))
	mux.Handle("GET /api/movie/{id}", errors.HttpErrorHandler(movieHandler.GetById))
	mux.Handle("POST /api/movie", errors.HttpErrorHandler(movieHandler.Create))

	mux.HandleFunc("GET /api/actors", actorHandler.GetAll)
	mux.HandleFunc("POST /api/actors", actorHandler.Create)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetByID)
	mux.HandleFunc("PATCH /api/actors/{id}", actorHandler.Update)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.Delete)
}
