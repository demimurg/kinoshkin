package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kinoshkin/domain"
)

type usersRepo struct {
	db *mongo.Database
}

func (u usersRepo) Get(id int) (*domain.User, error) {
	users := u.db.Collection("users")
	result := users.FindOne(ctx, bson.M{"_id": id})

	var dbUser bson.M
	err := result.Decode(&dbUser)
	if err != nil {
		return nil, err
	}

	var user domain.User
	user.ID, _ = dbUser["_id"].(int)
	user.Name, _ = dbUser["name"].(string)
	user.City, _ = dbUser["city"].(string)
	user.Long, user.Lat = extractLocation(dbUser)

	return &user, nil
}

func (u usersRepo) UpdateLoc(id int, lat, long float32) error {
	users := u.db.Collection("users")

	_, err := users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
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

func (u usersRepo) Create(id int, name string) error {
	users := u.db.Collection("users")

	_, err := users.InsertOne(ctx, bson.M{
		"_id":  id,
		"name": name,
		"city": "saint-petersburg",
		"location": bson.D{
			{"type", "Point"},
			{"coordinates", []float32{}},
		},
	})

	return err
}
