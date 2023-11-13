package main

import (
	"net/http"

	ps "github.com/brizaldi/go-parse"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var parser ps.Parser

	payload := ps.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = parser.WriteJSON(w, http.StatusOK, payload)
}
