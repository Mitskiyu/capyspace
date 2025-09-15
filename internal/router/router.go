package router

import (
	"database/sql"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	"github.com/Mitskiyu/capyspace/internal/database"
	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/space"
	"github.com/Mitskiyu/capyspace/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
)

func New(db *sql.DB, rdb *redis.Client, origins string) http.Handler {
	store := sqlc.New(db)
	cache := database.NewCache(rdb)

	authHandler := auth.NewHandler(auth.NewService(store, cache))
	userHandler := user.NewHandler(user.NewService(store))
	spaceHandler := space.NewHandler(space.NewService(store))

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{origins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Post("/users/check/email", userHandler.CheckEmail)
	r.Post("/users/check/username", userHandler.CheckUsername)
	r.Post("/users", authHandler.Register)
	r.Post("/sessions", authHandler.Login)

	r.Get("/spaces/{username}", spaceHandler.GetSpace)
	r.Group(func(r chi.Router) {
		r.Use(authHandler.SessionMiddleware)
		r.Post("/spaces", spaceHandler.CreateSpace)
	})

	return r
}
