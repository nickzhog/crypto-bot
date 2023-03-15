package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	DatabaseURI    string        `env:"DATABASE_URI"`
	UpdateInterval time.Duration `env:"POLL_INTERVAL"`
	TgToken        string        `env:"TELEGRAM_TOKEN"`
	ApiKey         string        `env:"API_KEY"`
}

func GetConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.DatabaseURI, "d", "", "postgres connection string")
	flag.DurationVar(&cfg.UpdateInterval, "i", time.Second*2, "interval for updates crypto")
	flag.StringVar(&cfg.TgToken, "a", "", "address for sending metrics")
	flag.StringVar(&cfg.ApiKey, "k", "", "api key")

	flag.Parse()

	env.Parse(&cfg)

	return cfg
}
