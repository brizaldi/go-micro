package main

import (
	"log-service/data"
	"net/http"

	ps "github.com/brizaldi/go-parse"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var parse ps.Parser

	// Read JSON into var
	var requestPayload JSONPayload
	_ = parse.ReadJSON(w, r, &requestPayload)

	// Insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		parse.ErrorJSON(w, err)
		return
	}

	resp := ps.JSONResponse{
		Error:   false,
		Message: "logged",
	}

	parse.WriteJSON(w, http.StatusAccepted, resp)
}
