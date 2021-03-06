package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kinoshkin/domain"
	"time"
)

func NewCinemasRepository(client mongo.Client) domain.CinemasRepo {
	return cinemasRepo{client.Database("kinoshkin")}
}

type cinemasRepo struct {
	db *mongo.Database
}

func (c cinemasRepo) Get(cinemaID string) (*domain.Cinema, error) {
	cinemas := c.db.Collection("cinemas")

	var cinema bson.M
	err := cinemas.FindOne(ctx, bson.M{"_id": cinemaID}).Decode(&cinema)
	if err != nil {
		return nil, err
	}
	return convertToDomainCinema(cinema), nil
}

func (c cinemasRepo) FindNearby(user *domain.User, pag domain.P) ([]*domain.Cinema, error) {
	geoNear := bson.D{
		{"$geoNear", bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": []float32{user.Long, user.Lat},
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

	cinemas := c.db.Collection("cinemas")
	docs, err := cinemas.Aggregate(ctx, mongo.Pipeline{
		geoNear,
		limit,
	})
	if err != nil {
		return nil, err
	}
	defer docs.Close(ctx)

	var cinemasList []*domain.Cinema
	for docs.Next(ctx) {
		var dbCinema bson.M
		if err := docs.Decode(&dbCinema); err != nil {
			return nil, err
		}

		cinemasList = append(cinemasList, convertToDomainCinema(dbCinema))
	}

	return cinemasList, nil
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

	movies := c.db.Collection("movies")
	docs, err := movies.Aggregate(ctx, mongo.Pipeline{
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
		var dbMovie bson.M
		if err := docs.Decode(&dbMovie); err != nil {
			return nil, err
		}

		moviesList = append(moviesList, domain.MovieWithSessions{
			Movie:    convertToDomainMovie(dbMovie),
			Sessions: extractSessions(dbMovie["schedule"]),
		})
	}

	return moviesList, nil
}
