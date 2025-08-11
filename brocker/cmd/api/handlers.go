package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *App) Brocker(wr http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Data successfuly received",
	}

	app.writeJson(wr, payload, http.StatusOK)
}

func (app *App) HandleSubmission(wr http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(wr, r, &requestPayload)
	if err != nil {
		app.errorJson(wr, err, http.StatusBadRequest)
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(wr, &requestPayload.Auth)
	default:
		app.errorJson(wr, fmt.Errorf("%s is not a valid action", requestPayload.Action), http.StatusBadRequest)
	}
}

func (app *App) authenticate(wr http.ResponseWriter, authPayload *AuthPayload) {
	// create json to send
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	// call the auth service
	request, reqErr := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if reqErr != nil {
		app.errorJson(wr, reqErr)
		return
	}

	client := &http.Client{}
	response, respErr := client.Do(request)
	if respErr != nil {
		app.errorJson(wr, respErr)
		return
	}
	defer response.Body.Close()

	// return correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(wr, fmt.Errorf("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(wr, fmt.Errorf("error calling auth service"), http.StatusUnauthorized)
		return
	}

	// on success answer
	var jsonFromService jsonResponse
	dErr := json.NewDecoder(response.Body).Decode(&jsonFromService)
	if dErr != nil {
		app.errorJson(wr, dErr)
		return
	}

	if jsonFromService.Error {
		app.errorJson(wr, fmt.Errorf(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	app.writeJson(wr, payload, http.StatusAccepted)
}
