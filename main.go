package main

import (
	"log"
	"net/http"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	sessionClient "github.com/good-threads/backend/internal/client/session"
	userClient "github.com/good-threads/backend/internal/client/user"
	"github.com/good-threads/backend/internal/config"
	commonLogic "github.com/good-threads/backend/internal/logic/common"
	sessionLogic "github.com/good-threads/backend/internal/logic/session"
	userLogic "github.com/good-threads/backend/internal/logic/user"
	httpPresentation "github.com/good-threads/backend/internal/presentation/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env := config.Get()
	mongoClient := mongoClient.Setup(env.MongoDBURI)
	userClient := userClient.Setup(mongoClient)
	httpPresentation := httpPresentation.Setup(
		commonLogic.Setup(),
		userLogic.Setup(userClient),
		sessionLogic.Setup(
			sessionClient.Setup(
				mongoClient,
			),
			userClient,
		),
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", httpPresentation.Ping)
	r.Post("/user", httpPresentation.CreateUser)
	r.Post("/session", httpPresentation.Login)
	r.Get("/board", httpPresentation.Ping)
	r.Patch("/board", httpPresentation.Ping)

	log.Println("Listening...")
	http.ListenAndServe(":3000", r)
}
