package service

import (
	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
)

type GenreService struct {
	repo repository.GenreRepository
}

func NewGenreService(repo repository.GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}
func (s *GenreService) CreateGenre(genre *entity.Genre) (int64, error) {
	if err := genre.Validate(); err != nil {
		return 0, err
	}
	id, err := s.repo.Create(genre)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (s *GenreService) GetAll() ([]entity.Genre, error) {
	genres, err := s.repo.GetAll()
	if err != nil {
		return []entity.Genre{}, err
	}
	return genres, nil
}
func (s *GenreService) GetByID(id int) (entity.Genre, error) {
	genre, err := s.repo.GetByID(id)
	if err != nil {
		return entity.Genre{}, err
	}
	return genre, nil
}
func (s *GenreService) Update(id int, genre entity.GenrePatchRequest) (entity.Genre, error) {
	genreUpdated, err := s.repo.Update(id, genre)
	if err != nil {
		return entity.Genre{}, err
	}
	return genreUpdated, nil
}
func (s *GenreService) Delete(id int, force bool) error {
	_, err := s.repo.Delete(id, force)
	if err != nil {
		return err
	}
	return nil
}
