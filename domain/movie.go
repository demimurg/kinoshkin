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
	// Operator the most underrated man
	Operator
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
	Description    string
	PosterURL      string
	Genre          Genre
	Duration       int
	AgeRestriction string
	FilmCrew       map[Position]Persons
	Rating         Rating
}

// CinemaWithSessions is a schedule of the certain movie in some cinema
type CinemaWithSessions struct {
	*Cinema
	Sessions []Session
}

// MoviesRepo work with movies collection
type MoviesRepo interface {
	Get(movID string) (*Movie, error)
	FindByRating(cityID string, pag P) ([]*Movie, error)
	GetSchedule(movieID, cityID string, pag P) ([]CinemaWithSessions, error)
}
