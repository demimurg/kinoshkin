package aggregator

import "go.mongodb.org/mongo-driver/mongo"

type Aggregator interface {
	Aggregate() error
}

func Cinemas(db *mongo.Database) Aggregator {
	return cinemaAgg{db: db}
}

func Cities(db *mongo.Database) Aggregator {
	return cityAgg{db: db}
}

func Schedule(db *mongo.Database) Aggregator {
	return &scheduleAgg{db: db, movies: make(map[string]*movie), emptyMovies: make(map[string]*movie)}
}
