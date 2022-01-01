package mongodb

import "kinoshkin/entity"

type user struct {
	ID       int    `bson:"_id"`
	Name     string `bson:"name"`
	City     string `bson:"city"`
	Location loc    `bson:"location"`
}

func toDomainUser(u *user) *entity.User {
	return &entity.User{
		ID:   u.ID,
		Name: u.Name,
		Lat:  u.Location.Coordinates.Latitude,
		Long: u.Location.Coordinates.Longitude,
		City: u.City,
	}
}
