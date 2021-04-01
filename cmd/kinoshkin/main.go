package main

import (
	"context"
	"kinoshkin/internal/bot"
	"kinoshkin/internal/conferencier"
	"kinoshkin/internal/mongodb"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.TODO()
	uri := "mongodb+srv://" + cfg.MongoUrl + "/kinoshkin?retryWrites=true&w=majority"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Can't connect to mongo database: ", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("kinoshkin")
	confSvc := conferencier.New(
		mongodb.NewCinemasRepository(db),
		mongodb.NewMoviesRepository(db),
		mongodb.NewUsersRepository(db),
		mongodb.NewSchedulesRepository(db),
	)

	bot.Start(confSvc)
}
