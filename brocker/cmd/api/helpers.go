package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
)


type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}


func (app *App) readJson(wr http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := math.Pow(2, 20) // one megabyte
	r.Body = http.MaxBytesReader(wr, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		log.Println("error occured inside readJosn: ", err)
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		err := errors.New("request Body must have only one JSON data")
		log.Println(err.Error())
		return err
	}
	
	return nil
}

func (app *App) writeJson(wr http.ResponseWriter, data any, status int, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	for key, val := range headers[0] {
		wr.Header()[key] = val
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(status)
	_, err = wr.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) errorJson(wr http.ResponseWriter, pErr error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error: true,
		Message: pErr.Error(),
	}

	wErr := app.writeJson(wr, payload, statusCode)
	if wErr != nil {
		return wErr
	}

	return nil
}