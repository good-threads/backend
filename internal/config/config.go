package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	TakenUsername string `env:"TAKEN_USERNAME,required"`
	MongoDBURI    string `env:"MONGO_DB_URI,required"`
}

func Get() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load config: %e", err)
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("unable parse config: %e", err)
	}
	return cfg
}
