package mongodb

import (
	"context"
	"kinoshkin/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMoviesRepository(db *mongo.Database) domain.MoviesRepository {
	return moviesRepo{db}
}

type moviesRepo struct {
	db *mongo.Database
}

var ctx = context.TODO()

func (m moviesRepo) Create(movs []domain.Movie) error {
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

func (m moviesRepo) Get(movID string) (*domain.Movie, error) {
	coll := m.db.Collection("movies")

	var mongoMov Movie
	err := coll.FindOne(ctx, bson.M{"_id": movID}).Decode(&mongoMov)
	if err != nil {
		return nil, err
	}

	mov := toDomainMovie(&mongoMov)
	return &mov, nil
}

func (m moviesRepo) FindByRating(city string, pag domain.P) ([]domain.Movie, error) {
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

	var movies []domain.Movie
	for docs.Next(ctx) {
		var mongoMov Movie
		if err := docs.Decode(&mongoMov); err != nil {
			return nil, err
		}
		movies = append(movies, toDomainMovie(&mongoMov))
	}

	return movies, nil
}

func (m moviesRepo) GetSchedule(
	movieID string, user *domain.User, pag domain.P,
) ([]domain.CinemaWithSessions, error) {
	calculateCinemasDist := bson.D{
		{"$geoNear", bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": []float64{user.Long, user.Lat},
			},
			"distanceField": "distance",
			"query": bson.M{
				"city_id": user.City,
			},
		}},
	}
	joinScheduleToCinemas := bson.D{
		{"$lookup", bson.M{
			"from": "schedule",
			"let":  bson.M{"id": "$_id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": bson.A{
					bson.M{"$eq": bson.A{"$cinema_id", "$$id"}},
					bson.M{"$eq": bson.A{"$movie_id", movieID}},
					bson.M{"$gte": bson.A{"$last", time.Now().Add(-10 * time.Minute)}},
					bson.M{"$lt": bson.A{"$last", time.Now().Truncate(24 * time.Hour).Add(27 * time.Hour)}},
				}}}},
			},
			"as": "schedule",
		}},
	}
	throwAwayCinemasWithEmptySchedule := bson.D{
		{"$match", bson.M{
			"schedule": bson.M{"$not": bson.M{"$size": 0}},
		}},
	}
	limit := bson.D{
		{"$limit", pag.Limit},
	}
	filterSchedule := bson.D{
		{"$addFields", bson.M{
			"schedule": bson.M{"$filter": bson.M{
				"input": bson.M{"$arrayElemAt": bson.A{"$schedule.showtimes", 0}},
				"as":    "session",
				"cond":  bson.M{"$gte": bson.A{"$$session.time", time.Now().Add(-10 * time.Minute)}},
			}},
		}},
	}

	docs, err := m.db.Collection("cinemas").Aggregate(ctx, mongo.Pipeline{
		calculateCinemasDist,
		joinScheduleToCinemas,
		throwAwayCinemasWithEmptySchedule,
		limit,
		filterSchedule,
	})
	if err != nil {
		return nil, err
	}
	defer docs.Close(ctx)

	var cinemasList []domain.CinemaWithSessions
	for docs.Next(ctx) {
		var movieSched movieSchedule
		if err := docs.Decode(&movieSched); err != nil {
			return nil, err
		}

		cinemasList = append(cinemasList, domain.CinemaWithSessions{
			Cinema:   toDomainCinema(&movieSched.Cinema),
			Sessions: toSessions(movieSched.Schedule),
		})
	}

	return cinemasList, nil
}
