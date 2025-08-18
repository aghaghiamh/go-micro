package main

import (
	"listener/internal/messagebroker"
	"listener/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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

	// Listen for messages

	// Create Consumer

	// Watch the queue and consume events
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
