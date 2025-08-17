package main

import (
	"fmt"
	"log"
	"mailer/domain"
	"mailer/mailer"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const webPort = "80"

type App struct {
	svc *mailer.Service
}

func main() {
	mailerSvc := mailer.New(createEmailServer())
	app := App{
		svc: &mailerSvc,
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

func createEmailServer() domain.EmailServer {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	return domain.EmailServer{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}
}
