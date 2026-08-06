package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kurbanamankeldi-alt/movies-api/customerrors"
	"github.com/kurbanamankeldi-alt/movies-api/entity"
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

func (h *MovieHandler) Get(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodGet {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	var (
		movies []*entity.Movie
		err    error
	)

	actorIdStr := r.URL.Query().Get("actor")
	genreIdStr := r.URL.Query().Get("genre")
	yearStr := r.URL.Query().Get("year")

	if actorIdStr != "" {
		actorIdInt, convertErr := strconv.Atoi(actorIdStr)
		if convertErr != nil || actorIdInt <= 0 {
			return &customerrors.HttpError{Message: "actor id should be positive number", Code: http.StatusBadRequest}
		}
		movies, err = h.service.FindMoviesByActor(actorIdInt)
	} else if genreIdStr != "" {
		genreIdInt, convertErr := strconv.Atoi(genreIdStr)
		if convertErr != nil || genreIdInt <= 0 {
			return &customerrors.HttpError{Message: "genre id should be positive number", Code: http.StatusBadRequest}
		}
		movies, err = h.service.FindMoviesByGenre(genreIdInt)
	} else if yearStr != "" {
		yearInt, convertErr := strconv.Atoi(yearStr)
		if convertErr != nil || yearInt <= 0 {
			return &customerrors.HttpError{Message: "release year should be positive number", Code: http.StatusBadRequest}
		}
		movies, err = h.service.FindMoviesByYear(yearInt)
	} else {
		movies, err = h.service.GetAllMovies()
	}

	if err != nil {
		return &customerrors.HttpError{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	response, err := json.Marshal(movies)

	if err != nil {
		return &customerrors.HttpError{Message: "failed to encode JSON", Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func (h *MovieHandler) GetById(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodGet {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		return &customerrors.HttpError{Message: "invalid movie id", Code: http.StatusBadRequest}
	}

	movie, err := h.service.GetMovieById(id)

	if err != nil {
		return &customerrors.HttpError{Message: "internal server error", Code: http.StatusInternalServerError}
	}

	response, err := json.Marshal(movie)

	if err != nil {
		return &customerrors.HttpError{Message: "failed to encode JSON", Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func (h *MovieHandler) GetActorsById(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodGet {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		return &customerrors.HttpError{Message: "invalid movie id", Code: http.StatusBadRequest}
	}

	actors, err := h.service.FindMovieActors(id)

	if err != nil {
		return &customerrors.HttpError{Message: "internal server error", Code: http.StatusInternalServerError}
	}

	response, err := json.Marshal(actors)

	if err != nil {
		return &customerrors.HttpError{Message: "failed to encode JSON", Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

// extra
func (h *MovieHandler) FilterBy(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodGet {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	idStr := r.URL.Query().Get("id")
	actorIdStr := r.URL.Query().Get("actorId")
	genreIdStr := r.URL.Query().Get("genreId")
	yearStr := r.URL.Query().Get("year")

	if idStr == "" && actorIdStr == "" && genreIdStr == "" && yearStr == "" {
		return h.Get(w, r)
	}

	var (
		id, actorId, genreId, year             int
		idErr, actorIdErr, genreIdErr, yearErr error
	)

	if idStr != "" {
		id, idErr = strconv.Atoi(idStr)
	}
	if actorIdStr != "" {
		actorId, actorIdErr = strconv.Atoi(actorIdStr)
	}
	if genreIdStr != "" {
		genreId, genreIdErr = strconv.Atoi(genreIdStr)
	}
	if yearStr != "" {
		year, yearErr = strconv.Atoi(yearStr)
	}

	if idErr != nil || actorIdErr != nil || genreIdErr != nil || yearErr != nil {
		return &customerrors.HttpError{Message: "param should be positive number", Code: http.StatusBadRequest}
	}

	if (idStr != "" && id <= 0) ||
		(actorIdStr != "" && actorId <= 0) ||
		(genreIdStr != "" && genreId <= 0) ||
		(yearStr != "" && year <= 0) {
		return &customerrors.HttpError{Message: "param should be positive number", Code: http.StatusBadRequest}
	}

	movie, err := h.service.FilterMoviesBy(id, actorId, genreId, year)

	if err != nil {
		return &customerrors.HttpError{Message: "internal server error", Code: http.StatusInternalServerError}
	}

	response, err := json.Marshal(movie)

	if err != nil {
		return &customerrors.HttpError{Message: "failed to encode JSON", Code: http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodPost {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	actorIdStr := r.URL.Query().Get("actors")
	genreIdStr := r.URL.Query().Get("genres")

	actorIds, errA := getIdsFromParam(actorIdStr)
	genreIds, errG := getIdsFromParam(genreIdStr)

	if errA != nil {
		return &customerrors.HttpError{Message: errA.Error(), Code: http.StatusBadRequest}
	}

	if errG != nil {
		return &customerrors.HttpError{Message: errG.Error(), Code: http.StatusBadRequest}
	}

	var movie entity.Movie

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&movie)
	if err != nil {
		return &customerrors.HttpError{Message: "invalid json", Code: http.StatusBadRequest}
	}

	for _, id := range actorIds {
		movie.Actors = append(movie.Actors, entity.Actor{Id: uint(id)})
	}

	for _, id := range genreIds {
		movie.Genres = append(movie.Genres, entity.Genre{Id: uint(id)})
	}

	createdId, err := h.service.CreateMovie(&movie)

	if err != nil {
		return &customerrors.HttpError{Message: "internal server error", Code: http.StatusInternalServerError}
	}

	movie.Id = uint(createdId)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)

	return nil
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodPatch {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	movieIdStr := r.PathValue("id")
	movieIdInt, err := strconv.Atoi(movieIdStr)

	if err != nil || movieIdInt <= 0 {
		return &customerrors.HttpError{Message: "invalid movie id", Code: http.StatusBadRequest}
	}

	var newData *entity.Movie

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	jsonErr := decoder.Decode(&newData)

	if jsonErr != nil {
		return &customerrors.HttpError{Message: jsonErr.Error(), Code: http.StatusBadRequest}
	}

	_, updateErr := h.service.UpdateMovie(movieIdInt, newData)

	if updateErr != nil {
		return &customerrors.HttpError{Message: updateErr.Error(), Code: http.StatusInternalServerError}
	}

	movie, err := h.service.GetMovieById(movieIdInt)

	if err != nil {
		return &customerrors.HttpError{Message: "internal server error", Code: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)

	return nil
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) *customerrors.HttpError {
	if r.Method != http.MethodDelete {
		return &customerrors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	movieIdStr := r.PathValue("id")
	movieIdInt, err := strconv.Atoi(movieIdStr)

	if err != nil || movieIdInt <= 0 {
		return &customerrors.HttpError{Message: "invalid movie id", Code: http.StatusBadRequest}
	}

	_, deleteErr := h.service.DeleteMovie(movieIdInt)

	if deleteErr != nil {
		return &customerrors.HttpError{Message: deleteErr.Error(), Code: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

// helper
func getIdsFromParam(param string) ([]int, error) {
	if len(param) == 0 {
		return []int{}, nil
	}

	param = strings.TrimSpace(param)
	paramArr := strings.Split(param, ",")
	ids := []int{}

	for _, val := range paramArr {
		id, err := strconv.Atoi(val)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid id provided: %v", id)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
