package aggregator

import (
	"kinoshkin/pkg/env"
)

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

var cfg = Config{}

func init() { env.Parse(&cfg) }
