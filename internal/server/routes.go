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

	successResponse(w, http.StatusOK, true)
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
		errorResponse(w, http.StatusInternalServerError, "Could not send email, try again later", err)
		return
	}

	successResponse(w, http.StatusOK, true)
}

func (s *Server) checkVerificationCodeHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	defer r.Body.Close()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("body decode error: %v", err))
		return
	}

	email := requestBody.Email
	pw := requestBody.Password

	if err := validate.Email(email); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}

	if err := validate.Password(pw); err != nil {
		errorResponse(w, http.StatusBadRequest, "Password is invalid", err)
		return
	}

	ctx := r.Context()
	_, err = auth.CreateUser(ctx, s.dbQueries, email, pw)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not sign up, try again later", err)
		return
	}

	successResponse(w, http.StatusOK, true)
}

func (s *Server) signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	defer r.Body.Close()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("body decode error: %v", err))
		return
	}

	email := requestBody.Email
	pw := requestBody.Password

	if err := validate.Email(email); err != nil {
		errorResponse(w, http.StatusBadRequest, "Email is invalid", err)
		return
	}

	if err := validate.Password(pw); err != nil {
		errorResponse(w, http.StatusBadRequest, "Password is invalid", err)
		return
	}

	ctx := r.Context()
	user, err := s.dbQueries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusUnauthorized, "No account associated with email", nil)
			return
		}
		errorResponse(w, http.StatusInternalServerError, "Could not log in, try again later", err)
		return
	}

	hashStr := ""
	if user.Password.Valid {
		hashStr = user.Password.String
	}

	saltStr := ""
	if user.Salt.Valid {
		saltStr = user.Salt.String
	}

	if hashStr == "" || saltStr == "" {
		errorResponse(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	match, err := auth.CheckPassword(pw, hashStr, saltStr)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Could not log in, try again later", err)
		return
	}

	if !match {
		errorResponse(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	sessionID, err := auth.CreateSession(ctx, s.dbQueries, user.ID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not log in, try again later", err)
		return
	}

	name := ""
	if user.Name.Valid {
		name = user.Name.String
	}

	token, err := auth.IssueSession(ctx, s.dbQueries, sessionID, user.ID, name, email, s.secretKey)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not log in, try again later", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30,
	})

	successResponse(w, http.StatusOK, true)
}
