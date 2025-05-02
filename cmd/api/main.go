package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/Mitskiyu/capyspace/internal/database"
	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/email"
	"github.com/Mitskiyu/capyspace/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load env file")
	}

	var (
		user           = os.Getenv("DB_USER")
		password       = os.Getenv("DB_PASSWORD")
		host           = os.Getenv("DB_HOST")
		port           = os.Getenv("DB_PORT")
		name           = os.Getenv("DB_NAME")
		secretKeyStr   = os.Getenv("JWT_SECRET_KEY")
		allowedOrigins = os.Getenv("CORS_ORIGINS")
	)

	if allowedOrigins == "" {
		log.Fatalf("Could not get CORS origins")
	}

	if secretKeyStr == "" {
		log.Fatalf("Could not get JWT secret key")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, name)

	dbConn := db.Connect(connStr)
	defer dbConn.Close()

	dbQueries := dbgen.New(dbConn)
	emailClient := email.New()
	secretKey := []byte(secretKeyStr)

	srv := server.New(dbConn, dbQueries, emailClient, secretKey, allowedOrigins)

	go func() {
		log.Printf("Capyspace server starting on port %s...", srv.Addr[1:])

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server forced to shutdown... %v", err)
	}
}
