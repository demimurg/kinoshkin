package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kinoshkin/internal/conferencier"
	"kinoshkin/internal/mongodb"
	bot "kinoshkin/internal/telebot"
	"log"
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

	bot.New(confSvc).Start()
}
