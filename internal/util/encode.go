package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func encode[T any](w http.ResponseWriter, status int, v T, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Println(err)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}
