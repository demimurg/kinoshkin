package domain

type User struct {
	ID        int
	Name      string
	Lat, Long float32
	City      string
}

type UsersRepo interface {
	Get(id int) (*User, error)
	// UpdateLoc used for updating user geo location
	UpdateLoc(id int, lat, long float32) error
	Create(id int, name string) error
}
