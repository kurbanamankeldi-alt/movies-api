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

func (s *MovieService) FindMoviesByGenre(id int) ([]*entity.Movie, error) {
	movies, err := s.repo.FindByGenre(id)

	if err != nil {
		return nil, err
	}

	return movies, nil	
}

func (s *MovieService) FindMoviesByYear(id int) ([]*entity.Movie, error) {
	movies, err := s.repo.FindByYear(id)

	if err != nil {
		return nil, err
	}

	return movies, nil	
}

func (s *MovieService) FindMovieActors(id int) ([]entity.Actor, error) {
	actors, err := s.repo.FindActors(id)

	if err != nil {
		return nil, err
	}

	return actors, nil	
}

func (s *MovieService) CreateMovie(movie *entity.Movie) (int64, error) {
	createdId, err := s.repo.Create(movie)

	if err != nil {
		return 0, err
	}

	return createdId, nil
}


func (s *MovieService) FilterMoviesBy(movieId, actorId, genreId, year int)  ([]*entity.Movie, error) {
	movies, err := s.repo.FilterBy(movieId, actorId, genreId, year)

	if err != nil {
		return nil, err
	}

	return movies, nil	
}

func (s *MovieService) UpdateMovie(id int, newData *entity.Movie) (int64, error) {
	updatedRow, err := s.repo.Update(id, newData)

	if err != nil {
		return 0, err
	}

	return updatedRow, nil
}