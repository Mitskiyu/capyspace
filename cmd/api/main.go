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

	"github.com/Mitskiyu/capyspace/internal/auth"
	db "github.com/Mitskiyu/capyspace/internal/database"
	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/email"
	"github.com/Mitskiyu/capyspace/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	// Pass the real os.Getenv function
	if err := run(ctx, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	// Create context that cancels on interrupt signal
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("could not load env file: %w", err)
	}

	var (
		env            = getenv("APP_ENV")
		srvPort        = getenv("SERVER_PORT")
		resendKey      = getenv("RESEND_API_KEY")
		secretKeyStr   = getenv("JWT_SECRET_KEY")
		allowedOrigins = getenv("CORS_ORIGINS")
		cookieDomain   = getenv("COOKIE_DOMAIN")
		dbUser         = getenv("DB_USER")
		dbPassword     = getenv("DB_PASSWORD")
		dbHost         = getenv("DB_HOST")
		dbPort         = getenv("DB_PORT")
		dbName         = getenv("DB_NAME")
	)

	if srvPort == "" {
		srvPort = "8080"
	}
	if resendKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable not set")
	}
	if allowedOrigins == "" {
		return fmt.Errorf("CORS_ORIGINS environment variable not set")
	}
	if secretKeyStr == "" {
		return fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}

	// Connect to db
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	dbConn, err := db.Connect(connStr)
	if err != nil {
		return fmt.Errorf("could not connect to db: %v", err)
	}
	defer dbConn.Close()

	dbq := dbgen.New(dbConn)
	emailClient := email.New(resendKey)
	sk := []byte(secretKeyStr)

	authHandler := &auth.Handler{
		DBQueries:   dbq,
		EmailClient: emailClient,
		SecretKey:   sk,
	}

	srv := server.Server{
		DBQueries:      dbq,
		SecretKey:      sk,
		AllowedOrigins: allowedOrigins,
		CookieDomain:   cookieDomain,
		CookieSecure:   env == "PROD",
	}

	srvHandler, err := server.New(authHandler, srv)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    ":" + srvPort,
		Handler: srvHandler,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Capyspace server starting on port %s...", srvPort)
		serverErrors <- httpServer.ListenAndServe()
	}()
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server failed: %w", err)
		}
	case <-ctx.Done():
		log.Println("Shutdown signal received, starting graceful shutdown...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		log.Println("Server gracefully stopped")
	}

	return nil
}
