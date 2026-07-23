package entity

type Movie struct {
	Id uint
	Title string
	ReleaseYear int
	Duration float64

	Actors []Actor
	Genres []Genre
}