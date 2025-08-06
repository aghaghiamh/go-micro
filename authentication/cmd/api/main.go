package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type App struct {
}

func main() {
	log.Println("Starting authentication service...")
	app := App{}

	srv := http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.SetRouter(),
	}

	log.Printf("The server is now running on %s Address.", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Because of the following error, server had to stopped: %s", err)
	}

}