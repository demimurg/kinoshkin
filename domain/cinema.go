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

// CinemasRepo provides basic db methods for cinemas collection
type CinemasRepo interface {
	// FindNearby search cinemas near user location
	FindNearby(lat, long float32) ([]*Cinema, error)
	FindWithMovie(lat, long float32, movieID string) ([]*Cinema, error)
	GetSchedule(cinemaID string) (map[*Movie][]*Session, error)
}
