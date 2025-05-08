package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccessBody struct {
	Data any `json:"data"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

func Success(w http.ResponseWriter, status int, data any) {
	res := SuccessBody{
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Could not encode success response: %v", err)
	}
}

func Error(w http.ResponseWriter, status int, msg string, err error) {
	res := ErrorBody{
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
