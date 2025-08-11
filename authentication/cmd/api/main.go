package main

import (
	"auth/data/adaptor"
	data "auth/data/repository"
	"auth/user"
	"fmt"
	"log"
	"net/http"
	"time"
)

const webPort = "80"

type App struct {
	Svc *user.Service
}

func main() {
	conn := adaptor.Retry(adaptor.OpenDB, 10, 2*time.Second)()
	userRepo := data.NewUserRepo(conn)

	userService := user.New(userRepo)
	log.Println("Starting authentication service...")
	app := App{
		Svc: userService,
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.SetRouter(),
	}

	log.Printf("The server is now running on %s Address.", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Because of the following error, server had to stopped: %s", err)
	}

}
