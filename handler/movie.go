package handler

import (
	"encoding/json"
	"net/http"

	//"strings"
	//"fmt"
	"strconv"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/errors"
	"github.com/kurbanamankeldi-alt/movies-api/service"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{
		service: s,
	}
}

func (h *MovieHandler) Get(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	if r.Method != http.MethodGet {
		return &errors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	var (
		movies []*entity.Movie
		err    error
	)

	actorIdStr := r.URL.Query().Get("actor")

	if actorIdStr != "" {
		actorIdInt, convertErr := strconv.Atoi(actorIdStr)
		if convertErr != nil || actorIdInt <= 0 {
			return &errors.HttpError{Message: "id should be positive number", Code: http.StatusBadRequest}
		}
		movies, err = h.service.FindMoviesByActor(actorIdInt)
	} else {
		movies, err = h.service.GetAllMovies()
	}

	if err != nil {
		return &errors.HttpError{Message: "movie not found", Code: http.StatusNotFound}
	}

	response, err := json.Marshal(movies)

	if err != nil {
		return &errors.HttpError{Message: "failed to encode JSON", Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func (h *MovieHandler) GetById(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	if r.Method != http.MethodGet {
		return &errors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return &errors.HttpError{Message: "invalid movie id", Code: http.StatusBadRequest}
	}

	movie, err := h.service.GetMovieById(id)

	if err != nil {
		return &errors.HttpError{Message: "movie not found", Code: http.StatusNotFound}
	}

	response, err := json.Marshal(movie)

	if err != nil {
		return &errors.HttpError{Message: "failed to encode JSON", Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	if r.Method != http.MethodPost {
		return &errors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	var movie entity.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		return &errors.HttpError{Message: "invalid json", Code: http.StatusBadRequest}
	}

	err = h.service.CreateMovie(&movie)

	if err != nil {
		return &errors.HttpError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)

	return nil
}
