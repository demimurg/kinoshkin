package aggregator

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var cfg = struct {
	// CinemasUrl is an api url for retrieving cinemas
	CinemasURL string `env:"CINEMAS_API_URL,required"`
	// CitiesURL is an api url for retrieving cities
	CitiesURL string `env:"CITIES_API_URL,required"`
	// ScheduleURL is an api url for retrieving schedule
	ScheduleURL string `env:"SCHEDULE_API_URL,required"`
	// TokenKP is a token for a movies service api
	TokenKP string `env:"KP_API_TOKEN,required"`
}{}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
}
