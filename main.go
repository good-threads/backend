package main

import (
	"log"
	"net/http"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	userClient "github.com/good-threads/backend/internal/client/user"
	"github.com/good-threads/backend/internal/config"
	commonLogic "github.com/good-threads/backend/internal/logic/common"
	userLogic "github.com/good-threads/backend/internal/logic/user"
	httpPresentation "github.com/good-threads/backend/internal/presentation/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env := config.Get()
	mongoClient := mongoClient.Setup(env.MongoDBURI)
	httpPresentation := httpPresentation.Setup(
		commonLogic.Setup(),
		userLogic.Setup(
			userClient.Setup(
				mongoClient,
			),
		),
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", httpPresentation.Ping)
	r.Post("/user", httpPresentation.CreateUser)

	log.Println("Listening...")
	http.ListenAndServe(":3000", r)
}
