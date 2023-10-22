package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"log"
)

type config struct {
	Token         string `env:"TELEGRAM_TOKEN"`
	DatabaseDSN   string `env:"DATABASE_DSN"'`
	TagBufferSize int
}

var cfg *config

func NewConfig() *config {
	if cfg != nil {
		return cfg
	}

	cfg = &config{}
	if err := env.Parse(&cfg); err != nil {
		cfg.Token = *flag.String("t", "", "token for telegram bot")
		cfg.DatabaseDSN = *flag.String("d", "user=postgres password=123456 host=localhost"+
			" port=5432 dbname=telegram", "database DSN format: user=user password=pass host=host port=port dbname=name")
		flag.Parse()
	}
	cfg.TagBufferSize = 20
	if cfg.Token == "" {
		log.Fatal("Empty token")
	}
	return cfg
}
