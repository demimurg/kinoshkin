package main

import (
	"context"
	"fmt"
	"kinoshkin/internal/aggregator"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var ctx = context.TODO()
	uri := "mongodb+srv://" + cfg.MongoAggUrl + "/kinoshkin?retryWrites=true&w=majority"
	mongodb, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Can't connect to mongo database: ", err)
	}
	defer mongodb.Disconnect(ctx)
	db := mongodb.Database("kinoshkin")

	for _, name := range cfg.Collections {
		var agg aggregator.Aggregator
		switch name {
		case "cities":
			agg = aggregator.Cities(db)
		case "cinemas":
			agg = aggregator.Cinemas(db)
		case "schedule":
			agg = aggregator.Schedule(db)
		default:
			fmt.Printf("wrong collection - %q, use cities/cinemas/movies", name)
		}

		if err := agg.Aggregate(); err != nil {
			log.Fatalf("%s aggregation failed: %s", name, err)
		}
	}
}
