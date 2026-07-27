package service

import (
	"github.com/kurbanamankeldi-alt/movies-api/entity"
	"github.com/kurbanamankeldi-alt/movies-api/repository"
)

type ActorService struct {
	repo repository.ActorRepository
}

func NewActorService(repo repository.ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}
func (s *ActorService) GetAll() ([]entity.Actor, error) {
	actors, err := s.repo.GetAll()
	if err != nil {
		return []entity.Actor{}, err
	}
	return actors, nil
}
