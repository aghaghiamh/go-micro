package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"brocker/internal/controller"
	"brocker/internal/events"
	rpcClient "brocker/internal/client"
	messagebroker "brocker/internal/message_broker"
	"brocker/utils"

	"github.com/joho/godotenv"
)

const webPort = 80

type Config struct {
	rabbitmq  messagebroker.Config
	backoff   utils.BackoffConfig
	publisher messagebroker.PublisherConfig
}

func main() {
	conf := loadConfig()

	// Connect to RabbitMQ
	client := messagebroker.NewClient(conf.rabbitmq)
	if err := utils.Backoff(client.Connect, conf.backoff)(); err != nil {
		log.Fatal("Couldn't connect to the message queue: ", err)
	}
	defer client.Close()

	// Rabbitmq Publisher
	publisher, pErr := messagebroker.NewPublisher(client, conf.publisher)
	if pErr != nil {
		log.Fatal("Couldn't create message broker publisher", pErr)
	}
	if err := publisher.Setup(conf.publisher); err != nil {
		log.Fatal("Couldn't set up exchange and channel for publisher: ", err)
	}

	eventSvc := events.NewEventService(publisher, "broker-service")

	// RPC Logger
	loggerRPCClient, loggerRPCErr := rpcClient.NewLoggerRPCClient("logger-service:5001")
	if loggerRPCErr != nil {
		// TODO: find a better approach than just failure on one service conenction
		log.Fatal("Logger RPC client service is not responsive: ", loggerRPCErr)
	}

	app := controller.NewApp(client, *eventSvc, loggerRPCClient)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.SetRouter(),
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

	// maxRetries, _ := strconv.Atoi(os.Getenv("BACKOFF_RABBIT_MAX_RETRIES"))

	return Config{
		rabbitmq: messagebroker.Config{
			Username: os.Getenv("RABBIT_USERNAME"),
			Password: os.Getenv("RABBIT_PASSWORD"),
			Host:     os.Getenv("RABBIT_HOST"),
			Port:     os.Getenv("RABBIT_PORT"),
		},
		backoff: utils.DefaultBackoffConfig(),
		// TODO: the DefaultPublisher should only be used for testing purpose; for full functionality should be override to get env values.
		publisher: messagebroker.DefaulPublisherConfig(),
	}
}
