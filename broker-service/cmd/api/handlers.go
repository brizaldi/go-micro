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
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
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

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	var parser ps.Parser

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		parser.ErrorJSON(w, err)
		return
	}

	payload := ps.JSONResponse{
		Error:   false,
		Message: "Logged",
	}

	parser.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	var parser ps.Parser

	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// Call the mail service
	mailServiceURL := "http://mail-service/send"

	// Post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		parser.ErrorJSON(w, errors.New("error calling mail service"))
		return
	}

	// Send back JSON
	payload := ps.JSONResponse{
		Error:   false,
		Message: "Message sent to " + msg.To,
	}

	parser.WriteJSON(w, http.StatusAccepted, payload)
}
