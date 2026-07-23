package handler

import(
	"encoding/json"
	"net/http"
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

func(h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var movie entity.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}	

	err = h.service.CreateMovie(&movie)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(movie)	
}