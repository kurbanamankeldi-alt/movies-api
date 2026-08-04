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
func (s *ActorService) CreateActor(actor *entity.Actor) (int64, error) {
	id, err := s.repo.Create(actor)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (s *ActorService) GetByID(id int) (entity.Actor, error) {
	actor, err := s.repo.GetByID(id)
	if err != nil {
		return entity.Actor{}, err
	}
	return actor, nil
}
func (s *ActorService) GetByName(name string) ([]entity.Actor, error) {
	actors, err := s.repo.GetByName(name)
	if err != nil {
		return []entity.Actor{}, err
	}
	return actors, nil
}
func (s *ActorService) Update(id int, actor entity.ActorPatchRequest) (entity.Actor, error) {
	actorUpdated, err := s.repo.Update(id, actor)
	if err != nil {
		return entity.Actor{}, err
	}
	return actorUpdated, nil
}
func (s *ActorService) Delete(id int, force bool) error {
	_, err := s.repo.Delete(id, force)
	if err != nil {
		return err
	}
	return nil
}
