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
		userId, expiring, err := h.service.sessionMiddleware(ctx, cookie.Value)
		if err != nil {
			log.Printf("%v at %s", err, r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}

		if expiring {
			refreshed := &http.Cookie{
				Name:     cookie.Name,
				Value:    cookie.Value,
				HttpOnly: cookie.HttpOnly,
				Secure:   cookie.Secure,
				SameSite: cookie.SameSite,
				MaxAge:   60 * 60 * 24 * 30, // 30 days
				Path:     cookie.Path,
			}

			http.SetCookie(w, refreshed)
		}

		ctx = context.WithValue(ctx, "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
