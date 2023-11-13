package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	ps "github.com/brizaldi/go-parse"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var parser ps.Parser

	payload := ps.JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = parser.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var parser ps.Parser
	var requestPayload RequestPayload

	err := parser.ReadJSON(w, r, &requestPayload)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		parser.ErrorJSON(w, errors.New("unknown action"))

	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	var parser ps.Parser

	// Create some JSON we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// Call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		parser.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		parser.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// Create a variable we'll read response.Body into
	var jsonFromService ps.JSONResponse

	// Decode the JSON from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		parser.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := ps.JSONResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	parser.WriteJSON(w, http.StatusAccepted, payload)
}
