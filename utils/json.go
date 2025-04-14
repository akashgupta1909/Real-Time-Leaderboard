package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError is a utility function to respond with an error message in JSON format.
func RespondWithError(responder http.ResponseWriter, code int, msg string) {
	if code >= 400 {
		log.Printf("Responding with %v code, msg: %v", code, msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJson(responder, code, errorResponse{
		Error: msg,
	})
}

// RespondWithJson is a utility function to respond with a JSON payload.
func RespondWithJson(responder http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error in json marshal for payload: %v", payload)
		responder.WriteHeader(500)
		return
	}
	responder.WriteHeader(code)
	responder.Header().Add("Content-Type", "application/json")
	responder.Write(data)
}
