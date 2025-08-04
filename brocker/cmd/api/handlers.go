package main

import (
	"encoding/json"
	"net/http"
)


type jsonResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}


func (app *App) Brocker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: "",
		Message: "Data successfuly received",
	}

	data, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}