package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"brocker/internal/controller"
	messagebroker "brocker/internal/message_broker"
	"brocker/utils"

	"github.com/joho/godotenv"
)

const webPort = 80

type Config struct {
	rabbitmq messagebroker.Config
	backoff  utils.BackoffConfig
}

func main() {
	conf := loadConfig()

	// Connect to RabbitMQ
	client := messagebroker.NewClient(conf.rabbitmq)
	err := utils.Backoff(client.Connect, conf.backoff)()
	if err != nil {
		log.Fatal("Couldn't connect to the message queue: ", err)
	}
	defer client.Close()

	app := controller.NewApp(client) 

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
	}
}
