package mongodb

import "kinoshkin/domain"

type user struct {
	ID       int    `bson:"_id"`
	Name     string `bson:"name"`
	City     string `bson:"city"`
	Location loc    `bson:"location"`
}

func toDomainUser(u *user) *domain.User {
	return &domain.User{
		ID:   u.ID,
		Name: u.Name,
		Lat:  u.Location.Coordinates.Latitude,
		Long: u.Location.Coordinates.Longitude,
		City: u.City,
	}
}
