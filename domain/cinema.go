package domain

import "time"

// Cinema represents movie theater
type Cinema struct {
	ID        string
	Name      string
	Address   string
	Metro     []string
	Lat, Long float32
	Distance  int
}

// Session shows information about session in cinema
type Session struct {
	ID    string
	Start time.Time
	Price int
}

// MovieWithSessions is a schedule for some movie in certain cinema
type MovieWithSessions struct {
	*Movie
	Sessions []Session
}

// CinemasRepo provides basic db methods for cinemas collection
type CinemasRepo interface {
	Get(cinemaID string) (*Cinema, error)
	// FindNearby search cinemas near user location
	FindNearby(lat, long float32, pag P) ([]*Cinema, error)
	FindWithMovie(lat, long float32, movieID string) ([]*Cinema, error)
	GetSchedule(cinemaID string, pag P) ([]MovieWithSessions, error)
}
