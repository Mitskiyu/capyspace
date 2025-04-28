package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/auth"
	"github.com/Mitskiyu/capyspace/internal/email"
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

	emailAddr := requestBody.Email
	if err := validate.Email(emailAddr); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}

	ctx := r.Context()
	res := make(map[string]bool)

	_, err = s.dbQueries.GetUserByEmail(ctx, emailAddr)
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

func (s *Server) sendVerificationHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx := r.Context()
	emailAddr := requestBody.Email
	if err := validate.Email(emailAddr); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}
	token, err := auth.CreateToken(ctx, s.dbQueries, emailAddr)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Service temporarily unavailable", err)
		return
	}

	emailConf := email.Email{
		To:       []string{emailAddr},
		From:     "noreply@capyspace.com",
		Subject:  "Capyspace Verification Code",
		HTMLBody: fmt.Sprintf("<h1>Code: %s</h1>", token),
		RawBody:  fmt.Sprintf("Code: %s", token),
	}

	if err := email.Send(ctx, s.emailClient, emailConf); err != nil {
		errorResponse(w, http.StatusBadRequest, "Could not send email, try again later", err)
		return
	}

	successResponse(w, http.StatusOK, "We sent a code to your inbox")
}
