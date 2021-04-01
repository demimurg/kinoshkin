package domain

type City struct {
	ID   string
	Name string
}

type CitiesRepository interface {
	Create(cities []City) error
}
