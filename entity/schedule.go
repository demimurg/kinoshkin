package entity

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
