package main

import "kinoshkin/pkg/env"

type Config struct {
	// MongoUrl is a mongodb cluster url for application
	MongoUrl string `env:"MONGO_URL,required"`
}

var cfg = Config{}

func init() { env.Parse(&cfg) }
