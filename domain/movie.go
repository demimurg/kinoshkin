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

// Rating for the most popular aggregators
type Rating struct {
	IMDB, KP float64
}

// Genre is a genre of the concrete movie
type Genre uint

const (
	Comedy Genre = iota + 1
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
	Duration       int32
	AgeRestriction string
	FilmCrew       map[Position][]string
	Rating         Rating
}

// CinemaWithSessions is a schedule of the certain movie in some cinema
type CinemaWithSessions struct {
	*Cinema
	Sessions []Session
}

// MoviesRepository work with movies collection
type MoviesRepository interface {
	Get(movID string) (*Movie, error)
	FindByRating(city string, pag P) ([]*Movie, error)
	GetSchedule(movieID string, user *User, pag P) ([]CinemaWithSessions, error)
}
