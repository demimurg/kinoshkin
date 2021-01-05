package env

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"sync"
)

var once sync.Once

// Parse wraps caarlos0/env package with .env loading
func Parse(v interface{}) {
	once.Do(func() {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Fatal("can't load variables from .env file")
		}
	})
	if err := env.Parse(v); err != nil {
		panic(errors.Wrap(err, "configuration setup failed"))
	}
}
