package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/jub0bs/cors"
)

type Server struct {
	dbConn      *sql.DB
	dbQueries   *dbgen.Queries
	emailClient *sesv2.Client
	secretKey   []byte
}

func New(dbConn *sql.DB, dbQueries *dbgen.Queries, emailClient *sesv2.Client, secretKey []byte, allowedOrigins string) *http.Server {
	s := &Server{
		dbConn:      dbConn,
		dbQueries:   dbQueries,
		emailClient: emailClient,
		secretKey:   secretKey,
	}

	corsMiddleware, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{allowedOrigins},
		Methods:        []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		RequestHeaders: []string{"Authorization", "Content-Type"},
	})
	if err != nil {
		log.Fatalf("CORS error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/auth/check-email", s.checkEmailHandler)
	mux.HandleFunc("/auth/send-verification", s.sendVerificationHandler)
	mux.HandleFunc("/auth/check-verification", s.checkVerficationCodeHandler)
	mux.HandleFunc("/auth/create-user", s.createUserHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware.Wrap(mux),
	}
}
