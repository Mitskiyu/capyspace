package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var hashPrms = hashParams{
	memory:  64 * 1024,
	iter:    3,
	threads: 2,
	saltLen: 16,
	keyLen:  32,
}

// Generates a random 8-digit code and returns it as a string
func generateCode() (string, error) {
	upper := big.NewInt(100000000) // Upper limit, 0 - 99999999

	n, err := rand.Int(rand.Reader, upper)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%08d", n.Int64()), nil
}

// Creates a verification code, stores it in verfication_codes
func createVerificationCode(ctx context.Context, dbq *dbgen.Queries, email string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}

	codeParams := dbgen.CreateVerificationCodeParams{
		ID:        uuid.New(),
		Email:     email,
		Code:      code,
		Used:      sql.NullTime{Valid: false},
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := dbq.CreateVerificationCode(ctx, codeParams); err != nil {
		return "", err
	}
	return code, nil
}

// Queries the db for a valid matching verification code and uses it
func checkVerificationCode(ctx context.Context, dbq *dbgen.Queries, email string, code string) (uuid.UUID, error) {
	codeParams := dbgen.GetValidVerificationCodeParams{
		Email: email,
		Code:  code,
	}

	// Code must match email, be unused and can't be expired
	id, err := dbq.GetValidVerificationCode(ctx, codeParams)
	if err != nil {
		return uuid.Nil, err
	}

	// Set the verification code to used
	if err := dbq.SetUsedVerificationCode(ctx, id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Creates a user and returns the user's id
// Checks if a user with the email already exists, email was verified and hashes the password using Argon2id
func createUser(ctx context.Context, dbq *dbgen.Queries, email string, pw string, code string) (uuid.UUID, error) {
	// Make sure no user with this email exists
	_, err := dbq.GetUserByEmail(ctx, email)
	if err == nil {
		return uuid.Nil, fmt.Errorf("user with email already exists")
	}
	if err != sql.ErrNoRows {
		return uuid.Nil, fmt.Errorf("query GetUserByEmail error: %w", err)
	}

	codeParams := dbgen.GetUsedVerificationCodeParams{
		Email: email,
		Code:  code,
	}

	// Check if the same code was used to verify
	_, err = dbq.GetUsedVerificationCode(ctx, codeParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("code is invalid or not used")
		}
		return uuid.Nil, fmt.Errorf("query GetUsedVerificationCode error: %w", err)
	}

	// Generate a random salt to randomize the hash
	salt := make([]byte, hashPrms.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return uuid.Nil, fmt.Errorf("could not generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(pw), salt, hashPrms.iter, hashPrms.memory, hashPrms.threads, hashPrms.keyLen)

	// Encodes hash and salt into strings to store in the db
	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)

	userParams := dbgen.CreateUserParams{
		ID:            uuid.New(),
		Name:          sql.NullString{Valid: false},
		Email:         email,
		Password:      sql.NullString{String: hashStr, Valid: true},
		Salt:          sql.NullString{String: saltStr, Valid: true},
		EmailVerified: time.Now(),
	}

	userID, err := dbq.CreateUser(ctx, userParams)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// Compares the password against the stored hash
func checkPassword(password, hashStr, saltStr string) (bool, error) {
	// Decodes the salt string to be used in the test hash
	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, fmt.Errorf("could not decode salt string: %w", err)
	}

	// Decodes the hash string to be used in comparison against the test hash
	hash, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return false, fmt.Errorf("could not decode hash string: %w", err)
	}

	test := argon2.IDKey([]byte(password), salt, hashPrms.iter, hashPrms.memory, hashPrms.threads, hashPrms.keyLen)

	// Compares the hash against the attempting password's test hash, and returns the result as a bool
	return bytes.Equal(hash, test), nil
}
