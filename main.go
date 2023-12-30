package main

import (
	"log"
	"net/http"

	changesetClient "github.com/good-threads/backend/internal/client/changeset"
	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	sessionClient "github.com/good-threads/backend/internal/client/session"
	threadClient "github.com/good-threads/backend/internal/client/thread"
	userClient "github.com/good-threads/backend/internal/client/user"
	"github.com/good-threads/backend/internal/config"
	boardLogic "github.com/good-threads/backend/internal/logic/board"
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
		boardLogic.Setup(
			userClient,
			changesetClient.Setup(
				mongoClient,
			),
			threadClient.Setup(
				mongoClient,
			),
		),
	)

	public := chi.NewRouter()
	public.Use(middleware.Logger)

	public.Get("/ping", httpPresentation.Ping)
	public.Post("/user", httpPresentation.CreateUser)
	public.Post("/session", httpPresentation.CreateSession)

	protected := public.Group(nil)
	protected.Use(httpPresentation.GetUsernameFromSession)

	protected.Get("/", httpPresentation.GetBoard)
	protected.Patch("/", httpPresentation.UpdateBoard)

	for _, s := range []string{
		"   ┓           ┓ ",
		"   ┃┏┏┓┏┓╋┏┓┏┓┏┫ ",
		"   ┛┗┛┗┗┛┗┣┛┗┻┗┻•",
		"          ┛      ",
		"    good threads.",
		"",
		"Listening...",
		"",
	} {
		log.Println(s)
	}
	http.ListenAndServe(":3000", public)
}
