package domain

import "time"

type User struct {
	ID        int
	Name      string
	Lat, Long float32
	City      string
	UpdatedAt time.Time
}

type UsersRepo interface {
	Get(id int) (*User, error)
	// UpdateLoc used for updating user geo location
	UpdateLoc(id int, lat, long float32) error
	GetLoc(id int) (lat float32, long float32, err error)
	GetCity(id int) (string, error)
	Create(id int, name string) error
}
