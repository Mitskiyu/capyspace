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

	"github.com/Mitskiyu/capyspace/internal/router"
)

func run() error {
	srv := http.Server{
		Addr:    ":80",
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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
