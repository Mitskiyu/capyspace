package router

import (
	"database/sql"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB) http.Handler {
	store := sqlc.New(db)
	authHandler := auth.NewHandler(auth.NewService(store))

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Post("/register", authHandler.Register)

	return r
}
