package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	TakenUsername string `env:"TAKEN_USERNAME,required"`
}

func Get() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
