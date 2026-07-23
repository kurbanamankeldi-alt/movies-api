package entity

import (
	"time"
)

type Actor struct {
	Id uint
	Name string
	BirthDate time.Time

	Movies []Movie
}