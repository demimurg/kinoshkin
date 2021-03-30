package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kinoshkin/domain"
	"time"
)

func NewCinemasRepository(db *mongo.Database) domain.CinemasRepository {
	return cinemasRepo{db}
}

type cinemasRepo struct {
	db *mongo.Database
}

func (c cinemasRepo) Create(cinemas []domain.Cinema) error {
	mongoCinemas := make([]interface{}, len(cinemas))
	for i := range cinemas {
		mongoCinemas[i] = toMongoCinema(&cinemas[i])
	}
	_, err := c.db.Collection("cinemas").InsertMany(ctx, mongoCinemas)
	return err
}

// todo: add distance information ($geoNear)
func (c cinemasRepo) Get(cinemaID string) (*domain.Cinema, error) {
	var cin Cinema
	err := c.db.Collection("cinemas").
		FindOne(ctx, bson.M{"_id": cinemaID}).Decode(&cin)
	if err != nil {
		return nil, err
	}
	domainCinema := toDomainCinema(&cin)

	return &domainCinema, nil
}

func (c cinemasRepo) FindNearby(user *domain.User, pag domain.P) ([]domain.Cinema, error) {
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

	var cinemas []domain.Cinema
	for docs.Next(ctx) {
		var mongoCinema Cinema
		if err := docs.Decode(&mongoCinema); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, toDomainCinema(&mongoCinema))
	}

	return cinemas, nil
}

func (c cinemasRepo) GetSchedule(cinemaID string, pag domain.P) ([]domain.MovieWithSessions, error) {
	joinScheduleToMovies := bson.D{
		{"$lookup", bson.M{
			"from": "schedule",
			"let":  bson.M{"id": "$_id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": bson.A{
					bson.M{"$eq": bson.A{"$cinema_id", cinemaID}},
					bson.M{"$eq": bson.A{"$movie_id", "$$id"}},
					bson.M{"$gte": bson.A{"$last", time.Now().Add(-10 * time.Minute)}},
					bson.M{"$lt": bson.A{"$last", time.Now().Truncate(24 * time.Hour).Add(27 * time.Hour)}},
				}}}},
			},
			"as": "schedule",
		}},
	}
	throwAwayMoviesWithEmptySchedule := bson.D{
		{"$match", bson.M{
			"schedule": bson.M{"$not": bson.M{"$size": 0}},
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
	filterSchedule := bson.D{
		{"$addFields", bson.M{
			"schedule": bson.M{"$filter": bson.M{
				"input": bson.M{"$arrayElemAt": bson.A{"$schedule.showtimes", 0}},
				"as":    "session",
				"cond":  bson.M{"$gte": bson.A{"$$session.time", time.Now().Add(-10 * time.Minute)}},
			}},
		}},
	}

	docs, err := c.db.Collection("movies").Aggregate(ctx, mongo.Pipeline{
		joinScheduleToMovies,
		throwAwayMoviesWithEmptySchedule,
		sortByRating,
		limit,
		filterSchedule,
	})
	if err != nil {
		return nil, err
	}
	defer docs.Close(ctx)

	var moviesList []domain.MovieWithSessions
	for docs.Next(ctx) {
		var cinSchedule cinemaSchedule
		if err := docs.Decode(&cinSchedule); err != nil {
			return nil, err
		}

		moviesList = append(moviesList, domain.MovieWithSessions{
			Movie:    toDomainMovie(&cinSchedule.Movie),
			Sessions: toSessions(cinSchedule.Schedule),
		})
	}

	return moviesList, nil
}
