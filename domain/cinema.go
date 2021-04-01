package domain

// Cinema represents movie theater
type Cinema struct {
	ID        string
	Name      string
	Address   string
	Metro     []string
	Lat, Long float64
	Distance  int
}

// CinemasRepository provides basic db methods for cinemas collection
type CinemasRepository interface {
	Create(cinemas []Cinema) error
	Get(cinemaID string) (*Cinema, error)
	// FindNearby search cinemas near user location
	FindNearby(user *User, pag P) ([]Cinema, error)
}
