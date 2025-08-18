package main

import (
	"auth/domain"
	"auth/user"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tsawler/toolbox"
)

type JsonResponse struct {
	Error   error  `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (app *App) Authenticate(wr http.ResponseWriter, r *http.Request) {
	var authRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var tools toolbox.Tools

	rjErr := tools.ReadJSON(wr, r, &authRequest)
	if rjErr != nil {
		log.Print("handler error: ", rjErr)
		tools.ErrorJSON(wr, rjErr, http.StatusBadRequest)

		return
	}

	_, svcErr := app.Svc.Authenticate(user.AuthRequest{
		Email:    authRequest.Email,
		Password: authRequest.Password,
	})
	if svcErr != nil {
		log.Print("handler error: ", svcErr)
		tools.ErrorJSON(wr, svcErr, http.StatusUnauthorized)

		return
	}

	logErr := app.writeLog(domain.LogEntry{
		Name: "Authentication",
		Data: fmt.Sprintf("%s logged in", authRequest.Email),
	})
	if logErr != nil {
		tools.ErrorJSON(wr, logErr)
	}

	tools.WriteJSON(wr, http.StatusAccepted, JsonResponse{
		Error:   nil,
		Message: "successfully authenticated!",
	})
}

func (app *App) writeLog(logEntry domain.LogEntry) error {
	jsonData, _ := json.MarshalIndent(logEntry, "", "\t")

	request, reqErr := http.NewRequest("POST", app.logSvcURL, bytes.NewBuffer(jsonData))
	if reqErr != nil {
		log.Println("Err: Couldn't create new request: ", reqErr, http.StatusInternalServerError)
		return reqErr
	}

	client := http.Client{}
	_, respErr := client.Do(request)
	if respErr != nil {
		log.Println("Err: Couldn't aggregate log inside log-service: ", respErr)
		return respErr
	}

	return nil
}
