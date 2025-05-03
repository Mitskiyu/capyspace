package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
)

type ctxKey string

const userCtxKey = ctxKey("user")

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			errorResponse(w, http.StatusUnauthorized, "Unauthorized request", fmt.Errorf("could not get session cookie: %v", err))
			return
		}

		token, claims, err := auth.RevalidateSession(r.Context(), s.dbQueries, s.secretKey, cookie.Value)
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, "Unauthorized request", err)
			return
		}

		if token != cookie.Value {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   60 * 60 * 24 * 30,
			})
		}

		ctx := context.WithValue(r.Context(), userCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
