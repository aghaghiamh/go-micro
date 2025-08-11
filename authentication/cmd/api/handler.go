package main

import (
	"auth/user"
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
		tools.ErrorJSON(wr, svcErr, http.StatusNotFound)

		return
	}

	tools.WriteJSON(wr, http.StatusAccepted, JsonResponse{
		Error:   nil,
		Message: "successfully authenticated!",
	})
}
