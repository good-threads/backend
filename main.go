package main

import (
	"context"
	"log"
	"net/http"

	userClient "github.com/good-threads/backend/internal/client/user"
	"github.com/good-threads/backend/internal/config"
	commonLogic "github.com/good-threads/backend/internal/logic/common"
	userLogic "github.com/good-threads/backend/internal/logic/user"
	httpPresentation "github.com/good-threads/backend/internal/presentation/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env, err := config.Get()
	if err != nil {
		log.Fatalf("unable get config: %e", err)
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.MongoDBURI))
	if err != nil {
		log.Fatalf("unable connect to mongo: %e", err)
	}

	httpPresentation := httpPresentation.Setup(
		commonLogic.Setup(),
		userLogic.Setup(
			userClient.Setup(mongoClient),
		),
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", httpPresentation.Ping)
	r.Post("/user", httpPresentation.CreateUser)

	log.Println("Listening...")
	http.ListenAndServe(":3000", r)
}
