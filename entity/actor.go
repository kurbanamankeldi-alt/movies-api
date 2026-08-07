package entity

import (
	"errors"
	"fmt"
	"time"
)

type Actor struct {
	Id        uint
	Name      string
	BirthDate time.Time

	Movies   []Movie `json:"movies,omitempty"`
	MovieIds []int   `json:"movieIds,omitempty"`
}

func (a *Actor) Validate() error {
	var err1, err2 error
	if a.Name == "" {
		err1 = fmt.Errorf("%w: name is empty", ErrInvalidContent)
	}
	if a.BirthDate.After(time.Now()) {
		err2 = fmt.Errorf("%w: the actor %s isn't born yet", ErrInvalidContent, a.Name)
	}
	return errors.Join(err1, err2)
}

type ActorPatchRequest struct {
	Name      *string
	BirthDate *string
	MovieIds  []int
}
