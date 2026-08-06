package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/errors"
	"github.com/kurbanamankeldi-alt/movies-api/service"
)

type ActorHandler struct {
	service *service.ActorService
}

func NewActorHandler(service *service.ActorService) *ActorHandler {
	return &ActorHandler{service: service}
}
func (h *ActorHandler) Create(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	var actor entity.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		return &errors.HttpError{Err: err, Message: "invalid json", Code: http.StatusBadRequest}
	}
	id, err := h.service.CreateActor(&actor)
	actor.Id = uint(id)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
	return nil
}
func (h *ActorHandler) GetAll(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	name := r.URL.Query().Get("name")
	gotMovies := r.URL.Query().Get("movies")
	movies := gotMovies == "true"
	if name == "" {
		actors, err := h.service.GetAll(movies)
		if err != nil {
			return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actors)
		return nil
	}
	actors, err := h.service.GetByName(name)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
	return nil
}
func (h *ActorHandler) GetByID(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return &errors.HttpError{Err: err, Message: "invalid id", Code: http.StatusBadRequest}
	}
	actor, err := h.service.GetByID(id)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
	return nil
}
func (h *ActorHandler) Update(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	idActor := r.PathValue("id")
	id, err := strconv.Atoi(idActor)
	if err != nil || id <= 0 {
		return &errors.HttpError{Err: err, Message: "invalid id", Code: http.StatusBadRequest}
	}
	var actorUpdate entity.ActorPatchRequest
	err1 := json.NewDecoder(r.Body).Decode(&actorUpdate)
	if err1 != nil {
		return &errors.HttpError{Err: err, Message: "invalid json", Code: http.StatusBadRequest}
	}
	actor, err := h.service.Update(id, actorUpdate)
	if err != nil {
		return &errors.HttpError{Err: err, Message: err.Error(), Code: http.StatusInternalServerError}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
	return nil
}
func (h *ActorHandler) Delete(w http.ResponseWriter, r *http.Request) *errors.HttpError {
	idActor := r.PathValue("id")
	id, err := strconv.Atoi(idActor)
	if err != nil || id <= 0 {
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
