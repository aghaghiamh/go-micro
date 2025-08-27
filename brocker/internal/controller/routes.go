package controller

import (
	"brocker/internal/client"
	"brocker/internal/events"
	messagebroker "brocker/internal/message_broker"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	rabbitmqClient *messagebroker.Client
	eventService   events.EventService
	Logger         *client.LoggerRPCClient
	LoggerGRPC     *client.LoggerGRPCClient
}

func NewApp(msgBrokerClient *messagebroker.Client, eventService events.EventService, logger *client.LoggerRPCClient, loggerGRPC *client.LoggerGRPCClient) App {
	return App{
		rabbitmqClient: msgBrokerClient,
		eventService:   eventService,
		Logger:         logger,
		LoggerGRPC:     loggerGRPC,
	}
}

func (app *App) SetRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Brocker)
	mux.Post("/handle", app.HandleSubmission)
	mux.Post("/logger-grpc", app.logGRPCItem)

	return mux
}
