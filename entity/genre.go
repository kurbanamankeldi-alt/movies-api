package entity

import "fmt"

type Genre struct {
	Id   uint
	Name string

	Movies   []Movie `json:"movies,omitempty"`
	MovieIds []int   `json:"movieIds,omitempty"`
}
type GenrePatchRequest struct {
	Name     *string
	MovieIds []int
}

func (g *Genre) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("%w: there is no name for genre", ErrInvalidContent)
	}
	return nil
}
