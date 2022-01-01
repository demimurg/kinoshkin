package mongodb

import (
	"kinoshkin/internal/entity"
	"kinoshkin/internal/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUsersRepository(db *mongo.Database) usecase.UsersRepository {
	return usersRepo{db}
}

type usersRepo struct {
	db *mongo.Database
}

func (u usersRepo) Get(id int) (*entity.User, error) {
	result := u.db.Collection("users").
		FindOne(ctx, bson.M{"_id": id})

	var mongoUser user
	err := result.Decode(&mongoUser)
	if err != nil {
		return nil, err
	}

	return toDomainUser(&mongoUser), nil
}

func (u usersRepo) UpdateLoc(id int, lat, long float32) error {
	// todo: use loc struct
	_, err := u.db.Collection("users").UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"location": bson.M{
			"type": "Point",
			"coordinates": bson.D{
				{"longitude", long},
				{"latitude", lat},
			},
		}},
	})

	return err
}

// todo: handle duplicate key error
func (u usersRepo) Create(id int, name string) error {
	_, err := u.db.Collection("users").InsertOne(ctx, user{
		ID:   id,
		Name: name,
		City: "saint-petersburg",
		Location: loc{
			Type: "Point",
			// Angleterre coordinates
			Coordinates: coords{
				Longitude: 30.308643285123765,
				Latitude:  59.93401788487186,
			},
		},
	})

	return err
}
