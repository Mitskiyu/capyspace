package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/jub0bs/cors"
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

	allowedOrigins := os.Getenv("CORS_ORIGINS")
	corsMiddleware, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{allowedOrigins},
		Methods:        []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		RequestHeaders: []string{"Authorization", "Content-Type"},
	})
	if err != nil {
		log.Fatalf("cors error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.healthHandler)
	mux.HandleFunc("/api/auth/check-email", s.checkEmailHandler)
	mux.HandleFunc("/api/auth/send-verification", s.sendVerificationHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware.Wrap(mux),
	}
}
