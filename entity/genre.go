package entity

type Genre struct {
	Id uint
	Name string

	Movies []Movie `json:"movies,omitempty"`
}