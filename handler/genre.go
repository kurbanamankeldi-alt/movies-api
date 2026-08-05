package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/errors"
	"github.com/kurbanamankeldi-alt/movies-api/service"
)

type GenreHandler struct {
	service *service.GenreService
}

func NewGenreHandler(service *service.GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

func (h *GenreHandler) Create(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	var genre entity.Genre
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		return &errors.HttpError{Err: err, Message: "invalid json", Code: http.StatusBadRequest}
	}
	id, err := h.service.CreateGenre(&genre)
	genre.Id = uint(id)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
	return nil
}
func (h *GenreHandler) GetAll(w http.ResponseWriter, r *http.Request) *errors.HttpError {

	genres, err := h.service.GetAll()
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
	return nil

}
func (h *GenreHandler) GetByID(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &errors.HttpError{Err: err, Message: "invalid id", Code: http.StatusBadRequest}
	}
	genre, err := h.service.GetByID(id)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
	return nil
}
func (h *GenreHandler) Update(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	idGenre := r.PathValue("id")
	id, err := strconv.Atoi(idGenre)
	if err != nil {
		return &errors.HttpError{Err: err, Message: "invalid id", Code: http.StatusBadRequest}
	}
	var genreUpdate entity.GenrePatchRequest
	err1 := json.NewDecoder(r.Body).Decode(&genreUpdate)
	if err1 != nil {
		return &errors.HttpError{Err: err, Message: "invalid json", Code: http.StatusBadRequest}
	}
	actor, err := h.service.Update(id, genreUpdate)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
	return nil
}
func (h *GenreHandler) Delete(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	idGenre := r.PathValue("id")
	id, err := strconv.Atoi(idGenre)
	if err != nil {
		return &errors.HttpError{Err: err, Message: "invalid id", Code: http.StatusBadRequest}
	}
	gotForce := r.URL.Query().Get("force")
	force := gotForce == "true"
	err = h.service.Delete(id, force)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
