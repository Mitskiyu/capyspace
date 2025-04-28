package main

import (
	"context"
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
)

func main() {
	dbConn := db.Connect()
	defer dbConn.Close()
	dbQueries := dbgen.New(dbConn)
	emailClient := email.New()

	srv := server.New(dbConn, dbQueries, emailClient)

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

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server forced to shutdown... %v", err)
	}
}
