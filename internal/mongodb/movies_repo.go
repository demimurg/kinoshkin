package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"kinoshkin/domain"
	"time"
)

func NewMoviesRepository(db *mongo.Database) domain.MoviesRepository {
	return moviesRepo{db}
}

type moviesRepo struct {
	db *mongo.Database
}

var ctx = context.TODO()

func (m moviesRepo) Get(movID string) (*domain.Movie, error) {
	coll := m.db.Collection("movies")

	var mov bson.M
	err := coll.FindOne(ctx, bson.M{"_id": movID}).Decode(&mov)
	if err != nil {
		return nil, err
	}
	return convertToDomainMovie(mov), nil
}

func (m moviesRepo) FindByRating(city string, pag domain.P) ([]*domain.Movie, error) {
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

	schedule := m.db.Collection("schedule")
	docs, err := schedule.Aggregate(ctx, mongo.Pipeline{
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

	var movies []*domain.Movie
	for docs.Next(ctx) {
		var dbMov bson.M
		if err := docs.Decode(&dbMov); err != nil {
			return nil, err
		}
		movies = append(movies, convertToDomainMovie(dbMov))
	}

	return movies, nil
}

func (m moviesRepo) GetSchedule(movieID string, user *domain.User, pag domain.P) ([]domain.CinemaWithSessions, error) {
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

	cinemas := m.db.Collection("cinemas")
	docs, err := cinemas.Aggregate(ctx, mongo.Pipeline{
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
		var dbCinema bson.M
		if err := docs.Decode(&dbCinema); err != nil {
			return nil, err
		}

		cinemasList = append(cinemasList, domain.CinemaWithSessions{
			Cinema:   convertToDomainCinema(dbCinema),
			Sessions: extractSessions(dbCinema["schedule"]),
		})
	}

	return cinemasList, nil
}

func extractSessions(showtimesB interface{}) []domain.Session {
	showtimes, ok := showtimesB.(bson.A)
	if !ok {
		return nil
	}

	var sessions []domain.Session
	for _, showI := range showtimes {
		show, ok := showI.(bson.M)
		if !ok {
			continue
		}

		var ses domain.Session
		ses.ID, _ = show["_id"].(string)
		ses.Price, _ = show["price"].(int32)

		bsonDate, ok := show["time"].(primitive.DateTime)
		if ok {
			ses.Start = bsonDate.Time()
		}

		sessions = append(sessions, ses)
	}

	return sessions
}

func convertToDomainCinema(dbCinema bson.M) *domain.Cinema {
	var cin domain.Cinema
	cin.ID, _ = dbCinema["_id"].(string)
	cin.Name, _ = dbCinema["name"].(string)
	cin.Address, _ = dbCinema["address"].(string)
	dist, ok := dbCinema["distance"].(float64)
	if ok {
		cin.Distance = int(dist)
	}
	metros, _ := dbCinema["metros"].(bson.A)
	for _, m := range metros {
		metro, ok := m.(string)
		if ok {
			cin.Metro = append(cin.Metro, metro)
		}
	}
	cin.Long, cin.Lat = extractLocation(dbCinema)

	return &cin
}

func extractLocation(doc bson.M) (long, lat float64) {
	location, ok := doc["location"].(bson.M)
	if ok {
		coordinates, _ := location["coordinates"].(bson.M)
		long, _ = coordinates["longitude"].(float64)
		lat, _ = coordinates["latitude"].(float64)
	}

	return
}

func convertToDomainMovie(dbMov bson.M) *domain.Movie {
	var mov domain.Movie
	mov.ID, _ = dbMov["_id"].(string)
	mov.Title, _ = dbMov["title"].(string)
	mov.Duration, _ = dbMov["duration"].(int32)
	mov.Description, _ = dbMov["description"].(string)
	mov.PosterURL, _ = dbMov["landscape_img"].(string)
	mov.AgeRestriction, _ = dbMov["age_restriction"].(string)

	rating, ok := dbMov["rating"].(bson.M)
	if ok {
		mov.Rating.KP, _ = rating["kp"].(float64)
		mov.Rating.IMDB, _ = rating["imdb"].(float64)

	}

	mov.FilmCrew = make(map[domain.Position][]string)
	staff, ok := dbMov["staff"].(bson.M)
	if ok {
		for role, personsI := range staff {
			personsI, ok := personsI.(bson.A)
			if !ok {
				continue
			}

			var persons []string
			for _, personI := range personsI {
				if p, ok := personI.(string); ok {
					persons = append(persons, p)
				}
			}

			switch role {
			case "actor":
				mov.FilmCrew[domain.Actor] = persons
			case "director":
				mov.FilmCrew[domain.Director] = persons
			case "writer":
				mov.FilmCrew[domain.Screenwriter] = persons
			case "operator":
				mov.FilmCrew[domain.Operator] = persons
			}
		}
	}

	return &mov
}
