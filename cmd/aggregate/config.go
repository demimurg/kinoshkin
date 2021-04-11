package main

import (
	"kinoshkin/pkg/env"
)

type Config struct {
	// Collections is subset of 'movies,cities,schedule'
	Collections []string `env:"COLLECTIONS" envDefault:"movies"`
	// MongoAggURL is a mongodb cluster url
	MongoAggURL string `env:"MONGO_AGG_URL,required"`
}

var cfg = Config{}

func init() { env.Parse(&cfg) }
