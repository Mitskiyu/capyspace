package space

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *handler) SpaceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIdRaw := ctx.Value("user_id")
		userID, ok := userIdRaw.(string)
		if !ok {
			log.Printf("mismatched type for user_id: %T", userIdRaw)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		spaceID := chi.URLParam(r, "spaceID")

		found, space, err := h.service.spaceMiddleware(ctx, userID)
		if err != nil {
			log.Printf("%v at %s", err, r.URL.Path)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !found {
			http.Error(w, "Space not found", http.StatusNotFound)
			return
		}

		if space.ID.String() != spaceID {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, "space", &space)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
