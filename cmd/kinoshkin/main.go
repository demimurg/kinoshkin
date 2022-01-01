package main

import (
	"context"
	"kinoshkin/internal/bot"
	"kinoshkin/internal/conferencier"
	"kinoshkin/internal/mongodb"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfg = struct {
	// mongodb cluster url for application
	MongoUrl string `env:"MONGO_URL,required"`
	// enable verbose logging for telebot
	BotToken string `env:"BOT_TOKEN,required"`
	// enable verbose logging for telebot
	VerboseLog bool `env:"VERBOSE_LOG" envDefault:"true"`
	// interval for polling telegram updates
	TelegramPoll time.Duration `env:"TELEGRAM_POLL" envDefault:"2s"`
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

	bot.Start(confSvc, cfg.BotToken, cfg.VerboseLog, cfg.TelegramPoll)
}
