package handler

import(
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/service"
	"fmt"
)

type MovieHandler struct {
    service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
    return &MovieHandler{
        service: s,
    }
}

func(h *MovieHandler) GetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	splitted := strings.Split(path, "/")
	fmt.Println(path, splitted)

	if len(splitted)<5 || splitted[2] != "movie" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	id := splitted[4]
	fmt.Println(id)
	convertedId, err := strconv.Atoi(id)


	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}	

	movie, err := h.service.GetMovieById(convertedId)

	if err != nil {
		http.Error(w, "movie not found", http.StatusNotFound)
		return		
	}

	response, err := json.Marshal(movie)

	if err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}	
	
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write(response)

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