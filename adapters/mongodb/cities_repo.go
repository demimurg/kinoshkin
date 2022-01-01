package mongodb

import (
	"kinoshkin/entity"
	"kinoshkin/usecase"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewCitiesRepository(db *mongo.Database) usecase.CitiesRepository {
	return citiesRepo{db: db}
}

type citiesRepo struct {
	db *mongo.Database
}

func (c citiesRepo) Create(cities []entity.City) error {
	mongoCities := make([]interface{}, len(cities))
	for i, domainCity := range cities {
		mongoCities[i] = city{
			ID:   domainCity.ID,
			Name: domainCity.Name,
		}
	}

	_, err := c.db.Collection("cities").InsertMany(ctx, mongoCities)
	return err
}

type city struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}
