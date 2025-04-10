package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
)

func New(db *sql.DB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbStatus := "connected"

		err := db.Ping()
		if err != nil {
			dbStatus = "disconnected: " + err.Error()
		}

		res := map[string]string{
			"status": status,
			"db":     dbStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}
