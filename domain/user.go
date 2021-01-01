package domain

import "time"

type User struct {
	ID        int
	Name      string
	Lat, Long float32
	UpdatedAt time.Time
}

type UsersRepo interface {
	// UpdateLoc used for updating user geo location
	UpdateLoc(id int, lat, long float32) error
	GetLoc(id int) (lat float32, long float32, err error)
	Create(id, name string) error
}
