package domain

// Position in team
type Position int

const (
	// Director of movie
	Director Position = iota
	// Screenwriter - person who wrote movie script
	Screenwriter
	// Actor (Angelina Jolie, Leonardo DiCaprio, etc)
	Actor
)

// Persons is a people working on one position
type Persons []string

// Rating for the most popular aggregators
type Rating struct {
	IMDB, KP float32
}

// Genre is a genre of the concrete movie
type Genre uint

const (
	Comedy Genre = iota
	Horror
	Action
)

// Movie is just a movie data
type Movie struct {
	ID             string
	Title          string
	PosterURL      string
	Genre          Genre
	Duration       int
	AgeRestriction bool
	FilmCrew       map[Position]Persons
	Rating         Rating
}

// MoviesRepo work with movies collection
type MoviesRepo interface {
	Find(cityID string) (*Movie, error)
	FindByRating() ([]*Movie, error)
	FindMany(ids ...string) ([]*Movie, error)
}
