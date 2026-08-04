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

	mux.Handle("GET /api/actors", errors.HttpErrorHandler(actorHandler.GetAll))
	mux.Handle("POST /api/actors", errors.HttpErrorHandler(actorHandler.Create))
	mux.Handle("GET /api/actors/{id}", errors.HttpErrorHandler(actorHandler.GetByID))
	mux.Handle("PATCH /api/actors/{id}", errors.HttpErrorHandler(actorHandler.Update))
	mux.Handle("DELETE /api/actors/{id}", errors.HttpErrorHandler(actorHandler.Delete))
}
