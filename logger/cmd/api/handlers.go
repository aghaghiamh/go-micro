package main

import (
	"log"
	"log-service/domain"
	"net/http"

	"github.com/tsawler/toolbox"
)

func (app *HTTPServer) WriteLog(wr http.ResponseWriter, r *http.Request) {
	var req WriteLogRequest
	var tools toolbox.Tools

	rjErr := tools.ReadJSON(wr, r, &req)
	if rjErr != nil {
		log.Print("handler error: ", rjErr)
		tools.ErrorJSON(wr, rjErr, http.StatusBadRequest)

		return
	}

	wlErr := app.svc.WriteLog(domain.LogEntry{
		Name: req.Name,
		Data: req.Data,
	})
	if wlErr != nil {
		log.Print("handler error: ", wlErr)
		tools.ErrorJSON(wr, rjErr, http.StatusInternalServerError)

		return
	}

	resp := toolbox.JSONResponse{
		Error:   false,
		Message: "logged",
	}
	tools.WriteJSON(wr, http.StatusAccepted, resp)
}
