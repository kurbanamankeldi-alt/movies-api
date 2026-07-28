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
	//genre
	genreRepo := repository.NewSQLiteGenreRepository(database)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handler.NewGenreHandler(genreService)

	mux.Handle("GET /api/movie", errors.HttpErrorHandler(movieHandler.Get))
	mux.Handle("GET /api/movie/{id}", errors.HttpErrorHandler(movieHandler.GetById))
	mux.Handle("GET /api/movie/filter", errors.HttpErrorHandler(movieHandler.FilterBy))
	mux.Handle("POST /api/movie", errors.HttpErrorHandler(movieHandler.Create))

	mux.Handle("GET /api/actors", errors.HttpErrorHandler(actorHandler.GetAll))
	mux.Handle("POST /api/actors", errors.HttpErrorHandler(actorHandler.Create))
	mux.Handle("GET /api/actors/{id}", errors.HttpErrorHandler(actorHandler.GetByID))
	mux.Handle("PATCH /api/actors/{id}", errors.HttpErrorHandler(actorHandler.Update))
	mux.Handle("DELETE /api/actors/{id}", errors.HttpErrorHandler(actorHandler.Delete))
	mux.Handle("DELETE /api/actors/deleteconnection/{id}", errors.HttpErrorHandler(actorHandler.DeleteConnection))

	mux.Handle("GET /api/genres", errors.HttpErrorHandler(genreHandler.GetAll))
	mux.Handle("POST /api/genres", errors.HttpErrorHandler(genreHandler.Create))
	mux.Handle("GET /api/genres/{id}", errors.HttpErrorHandler(genreHandler.GetByID))
	mux.Handle("PATCH /api/genres/{id}", errors.HttpErrorHandler(genreHandler.Update))
	mux.Handle("DELETE /api/genres/{id}", errors.HttpErrorHandler(genreHandler.Delete))
	mux.Handle("DELETE /api/genres/deleteconnection/{id}", errors.HttpErrorHandler(genreHandler.DeleteConnection))
}
