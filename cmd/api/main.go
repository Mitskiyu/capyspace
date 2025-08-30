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

	"github.com/Mitskiyu/capyspace/internal/database"
	"github.com/Mitskiyu/capyspace/internal/router"
	"github.com/Mitskiyu/capyspace/internal/util"
)

func run(getenv func(string, string) string) error {
	var (
		addr     = ":" + getenv("PORT", "8080")
		user     = getenv("DB_USER", "postgres")
		password = getenv("DB_PASSWORD", "postgres")
		host     = getenv("DB_HOST", "localhost")
		port     = getenv("DB_PORT", "5432")
		name     = getenv("DB_NAME", "capyspace")
	)

	log.Println("Connecting to postgres database...")
	db, err := database.Connect(user, password, host, port, name)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := database.Ping(db); err != nil {
		return err
	}
	log.Printf("Successfully connected to %s@%s:%s/%s", user, host, port, name)

	srv := http.Server{
		Addr:    addr,
		Handler: router.New(),
	}

	sig := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sig)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- fmt.Errorf("server failed to start: %v", err)
		}
	}()

	log.Printf("Server listening on %v...", srv.Addr)

	select {
	case <-sig:
		log.Println("Shutting down server...")
	case err := <-errs:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shut down: %v", err)
	}

	return nil
}

func main() {
	if err := run(util.GetEnv); err != nil {
		log.Fatal(err)
	}
}
