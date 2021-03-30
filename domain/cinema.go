package domain

import "time"

// Cinema represents movie theater
type Cinema struct {
	ID        string
	Name      string
	Address   string
	Metro     []string
	Lat, Long float64
	Distance  int
}

// Session shows information about session in cinema
type Session struct {
	ID    string
	Start time.Time
	Price int32
}

// MovieWithSessions is a schedule for some movie in certain cinema
type MovieWithSessions struct {
	Movie
	Sessions []Session
}

// CinemasRepository provides basic db methods for cinemas collection
type CinemasRepository interface {
	Create(cinemas []Cinema) error
	Get(cinemaID string) (*Cinema, error)
	// FindNearby search cinemas near user location
	FindNearby(user *User, pag P) ([]Cinema, error)
	GetSchedule(cinemaID string, pag P) ([]MovieWithSessions, error)
}
