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

	_, err = s.dbQueries.GetUserByEmail(ctx, emailAddr)
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, false)
		} else {
			errorResponse(w, http.StatusInternalServerError, "Could not check email, try again later", fmt.Errorf("database error: %v", err))
		}
		return
	}

	successResponse(w, http.StatusOK, false)
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

	emailAddr := requestBody.Email
	if err := validate.Email(emailAddr); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}

	ctx := r.Context()
	code, err := auth.CreateVerificationCode(ctx, s.dbQueries, emailAddr)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not send email, try again later", err)
		return
	}

	emailConf := email.Email{
		To:       []string{emailAddr},
		From:     "noreply@capyspace.com",
		Subject:  "Capyspace Verification Code",
		HTMLBody: fmt.Sprintf("<h1>Code: %s</h1>", code),
		RawBody:  fmt.Sprintf("Code: %s", code),
	}

	if err := email.Send(ctx, s.emailClient, emailConf); err != nil {
		errorResponse(w, http.StatusBadRequest, "Could not send email, try again later", err)
		return
	}

	successResponse(w, http.StatusOK, true)
}

func (s *Server) checkVerficationCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var requestBody struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	defer r.Body.Close()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("body decode error: %v", err))
		return
	}

	email := requestBody.Email
	code := requestBody.Code
	ctx := r.Context()

	if err := validate.Email(email); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}

	if err := validate.VerificationCode(code); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid or expired code", err)
		return
	}

	verified, err := auth.CheckVerificationCode(ctx, s.dbQueries, email, code)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not verify, try again later", err)
		return
	}

	if !verified {
		errorResponse(w, http.StatusBadRequest, "Invalid or expired code", nil)
		return
	}

	successResponse(w, http.StatusOK, true)
}
