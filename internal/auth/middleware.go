package auth

import (
	"context"
	"log"
	"net/http"
)

func (h *handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Printf("failed to get session cookie at %s: %v", r.URL.Path, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		userId, err := h.service.validateSession(ctx, cookie.Value)
		if err != nil {
			log.Printf("failed to validate session at %s: %v", r.URL.Path, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: call revalidate session if cache expiry date < 7 days

		ctx = context.WithValue(ctx, "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
