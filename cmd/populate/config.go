package main

import (
	"kinoshkin/pkg/env"
)

type Config struct {
	// Collections is subset of 'movies,cities,schedule'
	Collections []string `env:"COLLECTIONS" envDefault:"movies"`
	// MongoAggUrl is a mongodb cluster url
	MongoAggUrl string `env:"MONGO_AGG_URL,required"`
}

var cfg = Config{}

func init() { env.Parse(&cfg) }
