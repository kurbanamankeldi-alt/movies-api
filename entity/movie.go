package entity

type Movie struct {
	Id uint
	Title string
	ReleaseYear int
	Duration int

	Actors []Actor
	Genres []Genre
}