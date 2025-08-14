package main

import (
	"fmt"
	"log"
	"log-service/data/adaptor"
	"log-service/data/repository"
	"net/http"
	"os"

	"log-service/logSvc"

	"github.com/joho/godotenv"
)

const (
	webPort = "80"
)

type Config struct {
	mongo adaptor.MongoConfig
}

type App struct {
	svc *logsvc.Service
}

func main() {
	conf := loadConfig()
	client, err := adaptor.ConnectToMongo(conf.mongo)
	if err != nil {
		log.Panic(err)
	}
	// Close MongoDB connection
	defer adaptor.Disconnect(client)

	logRepo := repository.New(client)
	logSvc := logsvc.New(logRepo)

	app := App{
		svc: logSvc,
	}

	// start webserver
	app.serve()
}

func (app *App) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.SetRoutes(),
	}

	log.Printf("The server is now running on %s Address.", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Because of the following error, server had to stopped: %s", err)
	}
}

func loadConfig() Config {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	return Config{
		mongo: adaptor.MongoConfig{
			Username:     os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
			Password:     os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
			MongoAddress: os.Getenv("MONGO_ADDRESS"),
			DB:           os.Getenv("MONGO_INITDB_DATABASE"),
		},
	}
}
