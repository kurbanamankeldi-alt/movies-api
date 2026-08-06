package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"strconv"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/errors"
	"log"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{
		service: s,
	}
}

func(h *MovieHandler) Get(w http.ResponseWriter, r *http.Request) *errors.HttpError{
	log.Printf("Method: %q, Path: %q", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		return &errors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
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
			return &errors.HttpError{Message:"actor id should be positive number", Code: http.StatusBadRequest}	
		}
        movies, err = h.service.FindMoviesByActor(actorIdInt)
    } else if genreIdStr != "" {
		genreIdInt, convertErr := strconv.Atoi(genreIdStr)
		if convertErr != nil || genreIdInt <= 0 {
			return &errors.HttpError{Message:"genre id should be positive number", Code: http.StatusBadRequest}	
		}	
		movies, err = h.service.FindMoviesByGenre(genreIdInt)	
	} else if yearStr != "" {
		yearInt, convertErr := strconv.Atoi(yearStr)
		if convertErr != nil || yearInt <= 0 {
			return &errors.HttpError{Message:"release year should be positive number", Code: http.StatusBadRequest}	
		}	
		movies, err = h.service.FindMoviesByYear(yearInt)	
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

func(h *MovieHandler) GetActorsById(w http.ResponseWriter, r *http.Request) *errors.HttpError{
	if r.Method != http.MethodGet {
		return &errors.HttpError{Message:"method not allowed", Code: http.StatusMethodNotAllowed}
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return &errors.HttpError{Message:"invalid movie id", Code: http.StatusBadRequest}
	}	

	actors, err := h.service.FindMovieActors(id)

	if err != nil {
		return &errors.HttpError{Message:"no actors found", Code: http.StatusNotFound}		
	}

	response, err := json.Marshal(actors)

	if err != nil {
		return &errors.HttpError{Message:"failed to encode JSON", Code: http.StatusInternalServerError}		
	}	
	
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func(h *MovieHandler) FilterBy(w http.ResponseWriter, r *http.Request) *errors.HttpError{
	if r.Method != http.MethodGet {
		return &errors.HttpError{Message:"method not allowed", Code: http.StatusMethodNotAllowed}
	}

	idStr := r.URL.Query().Get("id")
	actorIdStr := r.URL.Query().Get("actorId")
	genreIdStr := r.URL.Query().Get("genreId")
	yearStr := r.URL.Query().Get("year")
	
	if idStr == "" && actorIdStr == "" && genreIdStr == "" &&  yearStr == "" {
		return h.Get(w,r)
	}

	var  (
		id, actorId, genreId, year int
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

	if idErr != nil || actorIdErr != nil || genreIdErr != nil || yearErr != nil{
		return &errors.HttpError{Message:"param should be positive number", Code: http.StatusBadRequest}
	}

	if (idStr != "" && id <=0) || 
	   (actorIdStr != "" && actorId <= 0) || 
	   (genreIdStr != "" && genreId <= 0) || 
	   (yearStr != "" && year <= 0) {
		return &errors.HttpError{Message:"param should be positive number", Code: http.StatusBadRequest}
	}


	movie, err := h.service.FilterMoviesBy(id, actorId, genreId, year)

	if err != nil {
		return &errors.HttpError{Message:"movie not found", Code: http.StatusNotFound}		
	}

	response, err := json.Marshal(movie)

	if err != nil {
		return &errors.HttpError{Message:"failed to encode JSON", Code: http.StatusInternalServerError}		
	}	
	
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil

}

func(h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) *errors.HttpError{
	if r.Method != http.MethodPost {
		return &errors.HttpError{Message: "method not allowed", Code: http.StatusMethodNotAllowed}
	}

	actorIdStr := r.URL.Query().Get("actors")
	genreIdStr := r.URL.Query().Get("genres")

	actorIds, errA := getIdsFromParam(actorIdStr)
	genreIds, errG := getIdsFromParam(genreIdStr)

	if errA != nil {
		return &errors.HttpError{Message:errA.Error(), Code: http.StatusBadRequest}
	}

	if errG != nil {
		return &errors.HttpError{Message:errG.Error(), Code: http.StatusBadRequest}
	}

	var movie entity.Movie

	for _, id := range actorIds {
		movie.Actors = append(movie.Actors, entity.Actor{Id:uint(id)})
	}

	for _, id := range genreIds {
		movie.Genres = append(movie.Genres, entity.Genre{Id:uint(id)})
	}	


	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		return &errors.HttpError{Message:"invalid json", Code: http.StatusBadRequest}
	}

	createdId, err := h.service.CreateMovie(&movie)

	if err != nil {
		return &errors.HttpError{Message: err.Error(), Code: http.StatusInternalServerError}
	}

	movie.Id = uint(createdId)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)

	return nil
}

//helper
func getIdsFromParam(param string) ([]int, error) {
	param = strings.TrimSpace(param)

	if len(param) == 0 {
		return nil, fmt.Errorf("params not provided: %v", param)
	}

	paramArr := strings.Split(param, ",")
	ids := []int{}

	for _, val := range paramArr {
		id, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid id provided: %v",  id)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
