package server

import (
	"fmt"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	"github.com/Mitskiyu/capyspace/internal/middleware"
	"github.com/Mitskiyu/capyspace/internal/response"
)

func handleRoutes(mux *http.ServeMux, handler *auth.Handler, authMw func(http.Handler) http.Handler) {
	mux.HandleFunc("/auth/check-email", handler.CheckEmail)
	mux.HandleFunc("/auth/send-verification", handler.SendVerification)
	mux.HandleFunc("/auth/check-verification", handler.CheckVerificationCode)
	mux.HandleFunc("/auth/create-user", handler.CreateUser)
	mux.HandleFunc("/auth/sign-in", handler.SignIn)

	mux.Handle("/mw/test", authMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.UserKey).(*auth.SessionClaims)
		if !ok {
			response.Error(w, http.StatusInternalServerError, "", fmt.Errorf("user claims not found in context"))
			return
		}
		response.Success(w, http.StatusOK, claims)
	})))
}
