package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/validate"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	successResponse(w, http.StatusOK, "OK")
}

func (s *Server) checkEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	defer r.Body.Close()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("email decode error: %v", err))
		return
	}

	email := requestBody.Email
	if err := validate.Email(email); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}

	ctx := r.Context()
	res := make(map[string]bool)

	_, err = s.dbQueries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			res["exists"] = false
			successResponse(w, http.StatusOK, res)
		} else {
			errorResponse(w, http.StatusInternalServerError, "Service temporarily unavailable", fmt.Errorf("database error: %v", err))
		}
		return
	}

	res["exists"] = true
	successResponse(w, http.StatusOK, res)
}
