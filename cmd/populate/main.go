package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kinoshkin/internal/aggregator"
	"log"
)

func main() {
	var ctx = context.TODO()
	mongodb, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://"+cfg.MongoAggUrl+"/kinoshkin?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal("Can't connect to mongo database: ", err)
	}
	defer mongodb.Disconnect(ctx)
	db := mongodb.Database("kinoshkin")

	for _, name := range cfg.Collections {
		var agg aggregator.Aggregator
		switch name {
		case "cities":
			// todo: remove later
			_ = db.Collection("cities").Drop(ctx)
			agg = aggregator.Cities(db)
		case "cinemas":
			_ = db.Collection("cinemas").Drop(ctx)
			agg = aggregator.Cinemas(db)
		case "movies":
			_ = db.Collection("tickets").Drop(ctx)
			_ = db.Collection("movies").Drop(ctx)
			agg = aggregator.Schedule(db)
		default:
			fmt.Printf("wrong collection - %q, use cities/cinemas/movies", name)
		}

		if err := agg.Aggregate(); err != nil {
			log.Fatalf("%s aggregation failed: %s", name, err)
		}
	}
}
