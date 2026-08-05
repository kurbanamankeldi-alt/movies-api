package entity

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
