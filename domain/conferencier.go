package domain

// P is a pagination info
type P struct {
	Offset, Limit int
}

// Conferencier is a main service, it manage all data about user, movie and cinema
type Conferencier interface {
	FindMovies(userID int, pagination P) ([]*Movie, error)
	FindCinemasNearby(userID int, pagination P) ([]*Cinema, []int, error)
	FindCinemasWithMovie(userID int, movieID string, pagination P) ([]*Cinema, []int, error)

	GetMovie(movieID string) (*Movie, error)
	GetCinema(cinemaID string) (*Cinema, map[*Movie][]Session, error)
	UpdateUserLocation(userID int, lat, long float32) error
	RegisterUser(userID int, name string) error
}
