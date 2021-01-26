package aggregator

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Aggregator interface {
	Aggregate() error
}

// todo: remove collection drop later
var ctx = context.TODO()

func Cinemas(db *mongo.Database) Aggregator {
	_ = db.Collection("cinemas").Drop(ctx)
	return cinemaAgg{db: db}
}

func Cities(db *mongo.Database) Aggregator {
	_ = db.Collection("cities").Drop(ctx)
	return cityAgg{db: db}
}

func Schedule(db *mongo.Database) Aggregator {
	return &scheduleAgg{db: db, movies: make(map[string]*movie), emptyMovies: make(map[string]*movie)}
}
