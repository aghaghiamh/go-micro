package controller

import (
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
}

func NewApp(msgBrokerClient *messagebroker.Client, eventService events.EventService) App {
	return App{
		rabbitmqClient: msgBrokerClient,
		eventService:   eventService,
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

	return mux
}
