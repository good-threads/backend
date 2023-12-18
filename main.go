package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Route string `env:"ROUTE,required"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get(cfg.Route, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome\n"))
	})

	log.Println("Listening...")

	http.ListenAndServe(":3000", r)
}
