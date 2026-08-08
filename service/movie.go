package service

import (
	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
	"strings"
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

func (s *MovieService) GetAllMovies(page, size int) ([]*entity.Movie, error) {
	movies, err := s.repo.FindAll(page, size)

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
	//make first letter upper and rest lower
	movie.Title = strings.Title(strings.ToLower(movie.Title))

	createdId, err := s.repo.Create(movie)

	if err != nil {
		return 0, err
	}

	return createdId, nil
}

func (s *MovieService) UpdateMovie(id int, newData *entity.Movie) (int64, error) {
	updatedRow, err := s.repo.Update(id, newData)

	if err != nil {
		return 0, err
	}

	return updatedRow, nil
}

func (s *MovieService) DeleteMovie(id int) (int64, error) {
	updatedRow, err := s.repo.Delete(id)

	if err != nil {
		return 0, err
	}

	return updatedRow, nil
}

//extra
func (s *MovieService) SearchMovies(title string) ([]*entity.Movie, error) {
	title = strings.ToLower(strings.TrimSpace(title))
	exactMatch := strings.Title(title)

	movie, err := s.repo.FindByExactTitle(exactMatch)
	if err != nil {
		return nil, err
	}

	if movie != nil {
		return []*entity.Movie{movie}, nil
	}

	return s.repo.FindByTitleContains(title)	
}