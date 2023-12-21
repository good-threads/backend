package main

import (
	"log"
	"net/http"

	"github.com/good-threads/backend/internal/config"
	"github.com/good-threads/backend/internal/handlers"
	"github.com/good-threads/backend/internal/logic/common"
	"github.com/good-threads/backend/internal/logic/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env, err := config.Setup()
	if err != nil {
		log.Fatalf("unable setup config: %e", err)
	}

	u := user.New()
	c := common.New()
	h := handlers.New(c, u)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get(env.Route, h.Ping)
	r.Post("/user", h.CreateUser)

	log.Println("Listening...")
	http.ListenAndServe(":3000", r)
}
