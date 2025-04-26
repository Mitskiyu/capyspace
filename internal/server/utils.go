package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func successResponse(w http.ResponseWriter, status int, data any) {
	res := SuccessResponse{
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Could not encode success response: %v", err)
	}
}

func errorResponse(w http.ResponseWriter, status int, msg string, err error) {
	res := ErrorResponse{
		Error: msg,
	}

	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Could not encode error response: %v", err)
	}
}
