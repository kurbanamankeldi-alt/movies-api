package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/service"
)

type ActorHandler struct {
	service *service.ActorService
}

func NewActorHandler(service *service.ActorService) *ActorHandler {
	return &ActorHandler{service: service}
}
func (h *ActorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		actors, err := h.service.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actors)
		return
	}
	actors, err := h.service.GetByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}
func (h *ActorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var actor entity.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	id, err := h.service.CreateActor(&actor)
	actor.Id = uint(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}
func (h *ActorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	actor, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}
func (h *ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	idActor := r.PathValue("id")
	id, err := strconv.Atoi(idActor)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var actorUpdate entity.ActorPatchRequest
	err1 := json.NewDecoder(r.Body).Decode(&actorUpdate)
	if err1 != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	actor, err := h.service.Update(id, actorUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}
