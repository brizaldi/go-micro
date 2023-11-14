package main

import (
	"net/http"

	ps "github.com/brizaldi/go-parse"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var parser ps.Parser

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := parser.ReadJSON(w, r, &requestPayload)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		parser.ErrorJSON(w, err)
		return
	}

	payload := ps.JSONResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	parser.WriteJSON(w, http.StatusAccepted, payload)

}
