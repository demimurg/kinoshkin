package aggregator

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
)

var cfg = Config{}

type Config struct {
	// CinemasUrl is an api url for retrieving cinemas
	CinemasURL string `env:"CINEMAS_API_URL,required"`
	// CitiesURL is an api url for retrieving cities
	CitiesURL string `env:"CITIES_API_URL,required"`
	// ScheduleURL is an api url for retrieving schedule
	ScheduleURL string `env:"SCHEDULE_API_URL,required"`
	// TokenKP is a token for a movies service api
	TokenKP string `env:"KP_API_TOKEN,required"`
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("can't load variables from .env file")
	}
	if err := env.Parse(&cfg); err != nil {
		panic(errors.Wrap(err, "configuration setup failed"))
	}
}
