package server

import (
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/middleware"

	"github.com/jub0bs/cors"
)

type Server struct {
	DBQueries      *dbgen.Queries
	SecretKey      []byte
	AllowedOrigins string
	CookieDomain   string
	CookieSecure   bool
}

func New(authHandler *auth.Handler, srv Server) (http.Handler, error) {
	corsMiddleware, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{srv.AllowedOrigins},
		Methods:        []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		RequestHeaders: []string{"Authorization", "Content-Type"},
		Credentialed:   true,
	})
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	authMw := middleware.NewAuthMiddleware(srv.DBQueries, srv.SecretKey, srv.CookieSecure, srv.CookieDomain)

	handleRoutes(mux, authHandler, authMw)

	return corsMiddleware.Wrap(mux), nil
}
