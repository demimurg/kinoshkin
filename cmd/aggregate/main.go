package main

import (
	"context"
	"fmt"
	"kinoshkin/internal/aggregator"
	"kinoshkin/internal/mongodb"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var ctx = context.TODO()
	uri := "mongodb+srv://" + cfg.MongoAggURL + "/kinoshkin?retryWrites=true&w=majority"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Can't connect to mongo database: ", err)
	}
	defer client.Disconnect(ctx)
	db := client.Database("kinoshkin")

	for _, name := range cfg.Collections {
		var agg aggregator.Aggregator
		switch name {
		case "cities":
			agg = aggregator.Cities(
				mongodb.NewCitiesRepository(db),
			)
		case "cinemas":
			agg = aggregator.Cinemas(
				mongodb.NewCinemasRepository(db),
			)
		case "schedule":
			agg = aggregator.Schedule(
				mongodb.NewMoviesRepository(db),
				mongodb.NewCinemasRepository(db),
				mongodb.NewSchedulesRepository(db),
			)
		default:
			fmt.Printf("wrong collection - %q, use cities/cinemas/schedule", name)
		}

		if err := agg.Aggregate(); err != nil {
			log.Printf("%s aggregation failed: %s", name, err)
		}
	}
}
