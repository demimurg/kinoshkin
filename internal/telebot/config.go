package telebot

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"time"
)

var cfg = Config{}

type Config struct {
	// Token is a telegram bot token
	Token string `env:"BOT_TOKEN,required"`
	// LogTrace sets trace log level for bot server
	LogTrace bool `env:"BOT_LOG_TRACE" envDefault:"false"`
	// UpdateInterval sets long pooling update interval for bot server
	UpdateInterval time.Duration `env:"BOT_UPDATE_INTERVAL" envDefault:"2s"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(errors.Wrap(err, "configuration setup failed"))
	}
}
