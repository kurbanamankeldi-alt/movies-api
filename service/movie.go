package service

import (
	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
)

type MovieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) CreateMovie(movie *entity.Movie) error {
	_, err := s.repo.Create(movie)

	if err != nil {
		return err
	}

	return nil
}

func (s *MovieService) GetMovieById(id int) (*entity.Movie, error) {
	movie, err := s.repo.FindById(id)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *MovieService) GetAllMovies() ([]*entity.Movie, error) {
	movies, err := s.repo.FindAll()

	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *MovieService) FindMoviesByActor(id int) ([]*entity.Movie, error) {
	movies, err := s.repo.FindByActor(id )

	if err != nil {
		return nil, err
	}

	return movies, nil	
}