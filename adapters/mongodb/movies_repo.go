package mongodb

import (
	"context"
	"kinoshkin/entity"
	"kinoshkin/usecase"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMoviesRepository(db *mongo.Database) usecase.MoviesRepository {
	return moviesRepo{db}
}

type moviesRepo struct {
	db *mongo.Database
}

var ctx = context.TODO()

func (m moviesRepo) Create(movs []entity.Movie) error {
	movies := m.db.Collection("movies")
	opts := &options.UpdateOptions{}
	opts = opts.SetUpsert(true)

	var mov interface{}
	for i := range movs {
		mov = toMongoMovie(&movs[i])
		_, err := movies.UpdateOne(ctx, bson.M{"_id": movs[i].ID}, bson.M{
			"$set": mov,
		}, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m moviesRepo) Get(movID string) (*entity.Movie, error) {
	coll := m.db.Collection("movies")

	var mongoMov Movie
	err := coll.FindOne(ctx, bson.M{"_id": movID}).Decode(&mongoMov)
	if err != nil {
		return nil, err
	}

	mov := toDomainMovie(&mongoMov)
	return &mov, nil
}

func (m moviesRepo) FindByRating(city string, pag usecase.P) ([]entity.Movie, error) {
	getFutureSessions := bson.D{
		{"$match", bson.D{
			{"city", city},
			{"last", bson.M{
				"$gte": time.Now().Add(-10 * time.Minute),
				"$lt":  time.Now().Truncate(24 * time.Hour).Add(27 * time.Hour),
			}},
		}},
	}
	groupByMovieID := bson.D{
		{"$group", bson.M{
			"_id": "$movie_id",
		}},
	}
	joinWithMovies := bson.D{
		{"$lookup", bson.M{
			"from":         "movies",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "movies",
		}},
	}
	throwAwayEmpty := bson.D{
		{"$match", bson.M{
			"movies": bson.M{"$not": bson.M{"$size": 0}},
		}},
	}
	extractMovieData := bson.D{
		{"$replaceRoot", bson.M{
			"newRoot": bson.M{
				"$arrayElemAt": bson.A{"$movies", 0},
			},
		}},
	}
	sortByRating := bson.D{
		{"$sort", bson.D{
			{"rating.kp", -1},
			{"rating.imdb", -1},
		}},
	}
	limit := bson.D{
		{"$limit", pag.Limit},
	}

	docs, err := m.db.Collection("schedule").Aggregate(ctx, mongo.Pipeline{
		getFutureSessions,
		groupByMovieID,
		joinWithMovies,
		throwAwayEmpty,
		extractMovieData,
		sortByRating,
		limit,
	})
	if err != nil {
		return nil, err
	}
	defer docs.Close(ctx)

	var movies []entity.Movie
	for docs.Next(ctx) {
		var mongoMov Movie
		if err := docs.Decode(&mongoMov); err != nil {
			return nil, err
		}
		movies = append(movies, toDomainMovie(&mongoMov))
	}

	return movies, nil
}
