package domain

import "time"

// Session shows information about session in cinema
type Session struct {
	ID    string
	Start time.Time
	Price int
}

// CinemaWithSessions is a schedule of the certain movie in some cinema
type CinemaWithSessions struct {
	Cinema
	Sessions []Session
}

// MovieWithSessions is a schedule for some movie in certain cinema
type MovieWithSessions struct {
	Movie
	Sessions []Session
}

type Schedule struct {
	MovieID  string
	CinemaID string
	Sessions []Session
}

type SchedulesRepository interface {
	Create(schedules []Schedule) error
	GetForMovie(movieID string, user *User, pag P) ([]CinemaWithSessions, error)
	GetForCinema(cinemaID string, pag P) ([]MovieWithSessions, error)
}
