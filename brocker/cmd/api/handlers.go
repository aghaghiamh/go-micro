package main

import (
	"net/http"
)

func (app *App) Brocker(wr http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Data successfuly received",
	}

	app.writeJson(wr, payload, http.StatusOK)
}
