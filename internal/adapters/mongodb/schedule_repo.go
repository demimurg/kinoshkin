package mongodb

import (
	"kinoshkin/internal/entity"
	"kinoshkin/internal/usecase"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewSchedulesRepository(db *mongo.Database) usecase.SchedulesRepository {
	return schedulesRepo{db}
}

type schedulesRepo struct {
	db *mongo.Database
}

func (s schedulesRepo) Create(schedules []entity.Schedule) error {
	mongoSchedules := make([]interface{}, len(schedules))
	for i, domainSched := range schedules {
		mongoSchedules[i] = toMongoSchedule(domainSched)
	}

	_, err := s.db.Collection("schedule").InsertMany(ctx, mongoSchedules)
	return err
}

func (s schedulesRepo) GetForMovie(movieID string, user *entity.User, pag usecase.P) ([]entity.CinemaWithSessions, error) {
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

	docs, err := s.db.Collection("cinemas").Aggregate(ctx, mongo.Pipeline{
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

	var cinemasList []entity.CinemaWithSessions
	for docs.Next(ctx) {
		var movieSched movieSchedule
		if err := docs.Decode(&movieSched); err != nil {
			return nil, err
		}

		cinemasList = append(cinemasList, entity.CinemaWithSessions{
			Cinema:   toDomainCinema(&movieSched.Cinema),
			Sessions: toDomainSessions(movieSched.Schedule),
		})
	}

	return cinemasList, nil
}

func (s schedulesRepo) GetForCinema(cinemaID string, pag usecase.P) ([]entity.MovieWithSessions, error) {
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

	docs, err := s.db.Collection("movies").Aggregate(ctx, mongo.Pipeline{
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

	var moviesList []entity.MovieWithSessions
	for docs.Next(ctx) {
		var cinSchedule cinemaSchedule
		if err := docs.Decode(&cinSchedule); err != nil {
			return nil, err
		}

		moviesList = append(moviesList, entity.MovieWithSessions{
			Movie:    toDomainMovie(&cinSchedule.Movie),
			Sessions: toDomainSessions(cinSchedule.Schedule),
		})
	}

	return moviesList, nil
}
