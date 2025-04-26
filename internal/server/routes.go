package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("email decode error: %v", err))
		return
	}

	if requestBody.Email == "" {
		errorResponse(w, http.StatusBadRequest, "Email is required", fmt.Errorf("email validation error: empty"))
		return
	}

	email := requestBody.Email
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
