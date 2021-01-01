package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

var cfg = Config{}

type Config struct {
	// Collections is subset of 'movies,cities,cinemas'
	Collections []string `env:"COLLECTIONS" envDefault:"movies"`
	// MongoAggUrl is a mongodb cluster url
	MongoAggUrl string `env:"MONGO_AGG_URL,required"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(errors.Wrap(err, "configuration setup failed"))
	}
}
