package main

import (
	"log"
	"mailer/domain"
	"net/http"

	"github.com/tsawler/toolbox"
)

func (app *App) SendMail(wr http.ResponseWriter, r *http.Request) {
	var req SendMailRequest
	var tools toolbox.Tools

	if rjErr := tools.ReadJSON(wr, r, &req); rjErr != nil {
		log.Println("Err: couldn't convert request into json: ", rjErr)
		tools.ErrorJSON(wr, rjErr, http.StatusBadRequest)

		return
	}

	msg := domain.Message{
		FromEmail: req.FromEmail,
		ToEmail:   req.ToEmail,
		Subject:   req.Subject,
		Data:      req.Message,
	}
	if smtpErr := app.svc.SendSMTPMessage(msg); smtpErr != nil {
		tools.ErrorJSON(wr, smtpErr, http.StatusBadRequest)

		return
	}

	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "sent to " + req.ToEmail,
	}

	tools.WriteJSON(wr, http.StatusOK, payload)
}
