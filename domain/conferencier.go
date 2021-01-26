package domain

// P is a pagination info
type P struct {
	Offset, Limit int
}

// Conferencier is a main service, it manage all data about user, movie and cinema
type Conferencier interface {
	FindMovies(userID int, pag P) ([]*Movie, error)
	FindCinemas(userID int, pag P) ([]*Cinema, error)

	GetMovie(movieID string) (*Movie, error)
	GetMovieSchedule(userID int, movieID string, pag P) ([]CinemaWithSessions, error)
	GetCinema(cinemaID string) (*Cinema, error)
	GetCinemaSchedule(cinemaID string, pag P) ([]MovieWithSessions, error)

	UpdateUserLocation(userID int, lat, long float32) error
	RegisterUser(userID int, name string) error
}
