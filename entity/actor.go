package entity

import (
	"time"
)

type Actor struct {
	Id        uint
	Name      string
	BirthDate time.Time

	Movies   []Movie `json:"movies,omitempty"`
	MovieIds []int   `json:"movieIds,omitempty"`
}
type ActorPatchRequest struct {
	Name      *string
	BirthDate *string
}
