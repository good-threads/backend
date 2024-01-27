package main

import (
	"log"
	"net/http"

	commandClient "github.com/good-threads/backend/internal/client/command"
	metricClient "github.com/good-threads/backend/internal/client/metric"
	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	sessionClient "github.com/good-threads/backend/internal/client/session"
	threadClient "github.com/good-threads/backend/internal/client/thread"
	userClient "github.com/good-threads/backend/internal/client/user"
	"github.com/good-threads/backend/internal/config"
	boardLogic "github.com/good-threads/backend/internal/logic/board"
	commonLogic "github.com/good-threads/backend/internal/logic/common"
	sessionLogic "github.com/good-threads/backend/internal/logic/session"
	threadLogic "github.com/good-threads/backend/internal/logic/thread"
	userLogic "github.com/good-threads/backend/internal/logic/user"
	httpPresentation "github.com/good-threads/backend/internal/presentation/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env := config.Get()
	metricClient := metricClient.Setup()
	mongoClient := mongoClient.Setup(env.MongoDBURI)
	userClient := userClient.Setup(mongoClient)
	threadClient := threadClient.Setup(mongoClient)
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
			commandClient.Setup(
				mongoClient,
			),
			threadClient,
			metricClient,
			mongoClient,
		),
		threadLogic.Setup(
			threadClient,
		),
	)

	public := chi.NewRouter()
	public.Use(middleware.Logger)
	public.Use(metricClient.Middleware)

	public.Get("/ping", httpPresentation.Ping)
	public.Handle("/metrics", metricClient.GetHandler())
	public.Post("/user", httpPresentation.CreateUser)
	public.Post("/session", httpPresentation.CreateSession)

	protected := public.Group(nil)
	protected.Use(httpPresentation.GetUsernameFromSession)

	protected.Get("/", httpPresentation.GetBoard)
	protected.Patch("/", httpPresentation.UpdateBoard)
	protected.Get("/thread/{id}", httpPresentation.GetThread)

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
