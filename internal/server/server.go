package server

import (
	"database/sql"
	"net/http"
	"os"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
)

type Server struct {
	dbConn    *sql.DB
	dbQueries *dbgen.Queries
}

func New(dbConn *sql.DB, dbQueries *dbgen.Queries) *http.Server {
	s := &Server{
		dbConn:    dbConn,
		dbQueries: dbQueries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.healthHandler)
	mux.HandleFunc("/api/auth/check-email", s.checkEmailHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}
