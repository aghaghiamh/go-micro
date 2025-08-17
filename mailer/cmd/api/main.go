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
	app := App{}

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
