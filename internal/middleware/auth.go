package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/response"
)

type CtxKey string

const UserKey CtxKey = "user"

// Maker function that revalidates JWT sessions
func NewAuthMiddleware(dbq *dbgen.Queries, sk []byte, secure bool, domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized request", nil)
				return
			}

			tokenStr := cookie.Value
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "Unauthorized request", nil)
				return
			}

			newTokenString, claims, err := auth.RevalidateSession(r.Context(), dbq, sk, tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized request", fmt.Errorf("AuthMiddleware RevalidateSession error: %v", err))
				return
			}
			if claims == nil {
				response.Error(w, http.StatusUnauthorized, "Unauthorized request", nil)
				return
			}

			if newTokenString != tokenStr {
				http.SetCookie(w, &http.Cookie{
					Name:     "session",
					Value:    newTokenString,
					Path:     "/",
					HttpOnly: true,
					Secure:   secure,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   60 * 60 * 24 * 90, // 90 days
					Domain:   domain,
				})
			}

			ctxWithUser := context.WithValue(r.Context(), UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		})
	}
}
