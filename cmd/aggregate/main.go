package main

import (
	"context"
	"fmt"
	"kinoshkin/internal/adapters/mongodb"
	"kinoshkin/internal/drivers/aggregator"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfg = struct {
	// subset of 'cities,cinemas,schedule' collections
	Collections []string `env:"COLLECTIONS" envDefault:"cities,cinemas,schedule"`
	// mongodb cluster url with auth params
	MongoURL string `env:"MONGO_URL,required"`
}{}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var ctx = context.TODO()
	uri := "mongodb+srv://" + cfg.MongoURL + "/kinoshkin?retryWrites=true&w=majority"
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
