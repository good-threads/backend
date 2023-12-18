package main

import (
	"log"
	"net/http"

	"github.com/good-threads/backend/internal/config"
	"github.com/good-threads/backend/internal/handlers"
	"github.com/good-threads/backend/internal/logic"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env, err := config.Setup()
	if err != nil {
		log.Fatalf("unable setup config: %e", err)
	}

	l := logic.NewTest()
	h := handlers.NewTest(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get(env.Route, h.Handler)

	log.Println("Listening...")
	http.ListenAndServe(":3000", r)
}
