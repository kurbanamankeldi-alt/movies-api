package server

import (
	"database/sql"
	"net/http"

	"github.com/kurbanamankeldi-alt/movies-api/handler"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
	"github.com/kurbanamankeldi-alt/movies-api/service"
)

func RegisterRoutes(mux *http.ServeMux, database *sql.DB) {
	movieRepo := repository.NewSQLiteMovieRepository(database)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handler.NewMovieHandler(movieService)

	mux.HandleFunc("/api/movie/", movieHandler.CreateMovie)
	mux.HandleFunc("/api/movie", movieHandler.CreateMovie)
}
