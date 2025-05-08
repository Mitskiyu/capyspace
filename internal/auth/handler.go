package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/email"
	res "github.com/Mitskiyu/capyspace/internal/response"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
)

// Handles http requests for auth
type Handler struct {
	DBQueries   *dbgen.Queries
	EmailClient *resend.Client
	SecretKey   []byte
}

// Checks if the provided email matches an existing user or has verified their email within an hour
func (h *Handler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var reqBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("/CheckEmail decode error: %v", err))
		return
	}

	email := reqBody.Email
	if err := validateEmail(email); err != nil {
		res.Error(w, http.StatusBadRequest, "Email is invalid", fmt.Errorf("/CheckEmail ValidateEmail error %v", err))
		return
	}

	ctx := r.Context()

	// Check if user exists
	_, err = h.DBQueries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Success(w, http.StatusOK, false)
			return
		}
		res.Error(w, http.StatusInternalServerError, "Could not check email, try again later", fmt.Errorf("/CheckEmail query GetUserByEmail error: %v", err))
		return
	}

	// User exists
	res.Success(w, http.StatusOK, true)
}

// Sends a verification code to the provided email
func (h *Handler) SendVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var reqBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("/SendVerification decode error: %v", err))
		return
	}

	emailAddr := reqBody.Email
	if err := validateEmail(emailAddr); err != nil {
		res.Error(w, http.StatusBadRequest, "Email is invalid", fmt.Errorf("/SendVerification validatEmail error: %v", err))
		return
	}

	ctx := r.Context()
	code, err := createVerificationCode(ctx, h.DBQueries, emailAddr)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Could not send email, try again later", fmt.Errorf("/SendVerification query createVerificationCode error: %v", err))
		return
	}

	emailParams := resend.SendEmailRequest{
		To:      []string{emailAddr},
		From:    "noreply@capyspace.com",
		Subject: "Capyspace Verification Code",
		Html:    fmt.Sprintf("<h1>Code: %s</h1>", code),
		Text:    fmt.Sprintf("Code: %s", code),
	}

	// Send email using email client
	_, err = email.Send(ctx, h.EmailClient, &emailParams)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Could not send email, try again later", fmt.Errorf("/SendVerification email.Send error: %v", err))
		return
	}

	// Email successfully sent
	res.Success(w, http.StatusOK, true)
}

// Checks that the verification code is valid for the email
func (h *Handler) CheckVerificationCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var reqBody struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("/CheckVerificationCode decode error: %v", err))
		return
	}

	email := reqBody.Email
	code := reqBody.Code
	ctx := r.Context()

	if err := validateEmail(email); err != nil {
		res.Error(w, http.StatusBadRequest, "Email is invalid", fmt.Errorf("/CheckVerificationCode validateEmail error: %v", err))
		return
	}

	if err := validateVerificationCode(code); err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid or expired code", fmt.Errorf("/CheckVerificationCode validateCode error: %v", err))
		return
	}

	// Check if verification code is valid
	id, err := checkVerificationCode(ctx, h.DBQueries, email, code)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Could not verify, try again later", fmt.Errorf("/CheckVerificationCode query checkVerificationCode error: %v", err))
		return
	}

	// Code is not valid
	if id == uuid.Nil {
		res.Success(w, http.StatusOK, false)
		return
	}

	res.Success(w, http.StatusOK, true)
}

// Creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("/CreateUser decode error: %v", err))
		return
	}

	email := reqBody.Email
	pw := reqBody.Password
	code := reqBody.Code

	if err := validateEmail(email); err != nil {
		res.Error(w, http.StatusBadRequest, "Email is invalid", fmt.Errorf("/CreateUser validateEmail error: %v", err))
		return
	}

	if err := validatePassword(pw); err != nil {
		res.Error(w, http.StatusBadRequest, "Password is invalid", fmt.Errorf("/CreateUser validatePassword error: %v", err))
		return
	}

	if err := validateVerificationCode(code); err != nil {
		res.Error(w, http.StatusBadRequest, "Code is invalid", fmt.Errorf("/CreateUser validateCode error: %v", err))
		return
	}

	ctx := r.Context()

	// Create user
	id, err := createUser(ctx, h.DBQueries, email, pw, code)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Could not sign up, try again later", fmt.Errorf("/CreateUser createUser error: %v", err))
		return
	}

	log.Printf("User with id: %v created", id)
	res.Success(w, http.StatusCreated, true)
}

// Authenticates a user and issues a JWT session
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()
	if err != nil {
		res.Error(w, http.StatusBadRequest, "Invalid request format", fmt.Errorf("/SignIn decode error: %v", err))
		return
	}

	email := reqBody.Email
	pw := reqBody.Password

	if err := validateEmail(email); err != nil {
		res.Error(w, http.StatusBadRequest, "Incorrect email or password", fmt.Errorf("/Signin validateEmail error: %v", err))
		return
	}

	if err := validatePassword(pw); err != nil {
		res.Error(w, http.StatusBadRequest, "Incorrect email or password", fmt.Errorf("/SignIn validatePassword error: %v", err))
		return
	}

	ctx := r.Context()

	// Check if user exists
	user, err := h.DBQueries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}
		res.Error(w, http.StatusInternalServerError, "Could not log in, try again later", err)
		return
	}

	hashStr := ""
	saltStr := ""

	if user.Password.Valid {
		hashStr = user.Password.String
	}

	if user.Salt.Valid {
		saltStr = user.Salt.String
	}

	if hashStr == "" || saltStr == "" {
		res.Error(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// Check if password matches
	match, err := checkPassword(pw, hashStr, saltStr)
	if err != nil {
		res.Error(w, http.StatusUnauthorized, "Could not log in, try again later", fmt.Errorf("/SignIn checkPassword error: %v", err))
		return
	}

	// Password doesn't match
	if !match {
		res.Error(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	// Create the db session
	sessionID, err := createSession(ctx, h.DBQueries, user.ID)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Could not log in, try again later", fmt.Errorf("/SignIn createSession error: %v", err))
		return
	}

	name := ""
	if user.Name.Valid {
		name = user.Name.String
	}

	// Issue JWT session
	token, err := issueSession(ctx, h.DBQueries, sessionID, user.ID, name, email, h.SecretKey)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "Could not log in, try again later", fmt.Errorf("/SignIn issueSession error: %v", err))
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

	// Sucessfully signed in
	res.Success(w, http.StatusOK, true)
}
