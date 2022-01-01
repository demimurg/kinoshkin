package mongodb

import (
	"kinoshkin/entity"
	"kinoshkin/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCinemasRepository(db *mongo.Database) usecase.CinemasRepository {
	return cinemasRepo{db}
}

type cinemasRepo struct {
	db *mongo.Database
}

func (c cinemasRepo) Create(cinemas []entity.Cinema) error {
	mongoCinemas := make([]interface{}, len(cinemas))
	for i := range cinemas {
		mongoCinemas[i] = toMongoCinema(&cinemas[i])
	}
	_, err := c.db.Collection("cinemas").InsertMany(ctx, mongoCinemas)
	return err
}

// todo: add distance information ($geoNear)
func (c cinemasRepo) Get(cinemaID string) (*entity.Cinema, error) {
	var cin Cinema
	err := c.db.Collection("cinemas").
		FindOne(ctx, bson.M{"_id": cinemaID}).Decode(&cin)
	if err != nil {
		return nil, err
	}
	domainCinema := toDomainCinema(&cin)

	return &domainCinema, nil
}

func (c cinemasRepo) GetAll(cityID string) ([]entity.Cinema, error) {
	cursor, err := c.db.Collection("cinemas").Find(ctx, bson.M{"city_id": cityID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cinemas []entity.Cinema
	for cursor.Next(ctx) {
		var mongoCinema Cinema
		if err := cursor.Decode(&mongoCinema); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, toDomainCinema(&mongoCinema))
	}

	return cinemas, nil
}

func (c cinemasRepo) FindNearby(user *entity.User, pag usecase.P) ([]entity.Cinema, error) {
	geoNear := bson.D{
		{"$geoNear", bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": []float64{user.Long, user.Lat},
			},
			"maxDistance":   20000, // in meters
			"distanceField": "distance",
			"query": bson.M{
				"city_id": user.City,
			},
		}},
	}
	limit := bson.D{
		{"$limit", pag.Limit},
	}

	docs, err := c.db.Collection("cinemas").Aggregate(ctx, mongo.Pipeline{
		geoNear,
		limit,
	})
	if err != nil {
		return nil, err
	}
	defer docs.Close(ctx)

	var cinemas []entity.Cinema
	for docs.Next(ctx) {
		var mongoCinema Cinema
		if err := docs.Decode(&mongoCinema); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, toDomainCinema(&mongoCinema))
	}

	return cinemas, nil
}
