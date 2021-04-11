package domain

import "time"

// Position in team
type Position = string

const (
	// Director of movie
	Director Position = "режиссер"
	// Screenwriter - person who wrote movie script
	Screenwriter Position = "сценарист"
	// Actor (Angelina Jolie, Leonardo DiCaprio, etc)
	Actor Position = "актер"
	// Operator the most underrated man
	Operator Position = "оператор"
)

// Rating for the most popular aggregators
type Rating struct {
	IMDB, KP float64
}

// Movie is just a movie data
type Movie struct {
	ID             string
	KpID           string
	Title          string
	Description    string
	PosterURL      string
	Duration       int
	AgeRestriction string
	FilmCrew       map[Position][]string
	Rating         Rating
	DateReleased   time.Time
}

// MoviesRepository work with movies collection
type MoviesRepository interface {
	Create(movs []Movie) error
	Get(movID string) (*Movie, error)
	FindByRating(city string, pag P) ([]Movie, error)
}
