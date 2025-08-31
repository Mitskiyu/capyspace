package router

import (
	"database/sql"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	"github.com/Mitskiyu/capyspace/internal/database"
	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func New(db *sql.DB, rdb *redis.Client) http.Handler {
	store := sqlc.New(db)
	cache := database.NewCache(rdb)

	authHandler := auth.NewHandler(auth.NewService(store, cache))

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Post("/register", authHandler.Register)

	return r
}
