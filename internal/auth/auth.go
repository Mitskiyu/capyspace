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

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type SessionClaims struct {
	SessionID    uuid.UUID `json:"session_id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RevalidateAt int64     `json:"revalidate_at"`
	jwt.RegisteredClaims
}

func GenerateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("could not generate code: %v", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func CreateVerificationCode(ctx context.Context, dbQueries *dbgen.Queries, email string) (string, error) {
	code, err := GenerateVerificationCode()
	if err != nil {
		return "", err
	}

	codeParams := dbgen.CreateVerificationCodeParams{
		ID:        uuid.New(),
		Email:     email,
		Code:      code,
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := dbQueries.CreateVerificationCode(ctx, codeParams); err != nil {
		return "", fmt.Errorf("could not create code: %v", err)
	}
	return code, nil
}

func CheckVerificationCode(ctx context.Context, dbQueries *dbgen.Queries, email string, code string) (bool, error) {
	codeParams := dbgen.GetValidVerificationCodeParams{
		Email: email,
		Code:  code,
	}

	id, err := dbQueries.GetValidVerificationCode(ctx, codeParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("could not get valid code: %v", err)
	}

	// set the code's <used> to true
	if err := dbQueries.SetUsedVerificationCode(ctx, id); err != nil {
		return false, fmt.Errorf("could not set used code: %v", err)
	}

	return true, nil
}

func CreateUser(ctx context.Context, dbQueries *dbgen.Queries, email string, pw string) (uuid.UUID, error) {
	_, err := dbQueries.GetUsedVerificationCode(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("could not create user: email not verified")
		}
		return uuid.Nil, fmt.Errorf("could not get used code: %v", err)
	}

	var (
		memory     uint32 = 64 * 1024
		iter       uint32 = 3
		threads    uint8  = 2
		saltLength uint32 = 16
		keyLen     uint32 = 32
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return uuid.Nil, fmt.Errorf("could not generate salt: %v", err)
	}

	hash := argon2.IDKey([]byte(pw), salt, iter, memory, threads, keyLen)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)
	saltStr := base64.RawStdEncoding.EncodeToString(salt)

	userParams := dbgen.CreateUserParams{
		ID:            uuid.New(),
		Name:          sql.NullString{String: "", Valid: false},
		Email:         email,
		Password:      sql.NullString{String: hashStr, Valid: true},
		Salt:          sql.NullString{String: saltStr, Valid: true},
		EmailVerified: time.Now(),
	}

	id, err := dbQueries.CreateUser(ctx, userParams)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not create user: %v", err)
	}

	return id, nil
}

func CreateSession(ctx context.Context, dbQueries *dbgen.Queries, userID uuid.UUID) (uuid.UUID, error) {
	sessionParams := dbgen.CreateSessionParams{
		ID:        uuid.New(),
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	id, err := dbQueries.CreateSession(ctx, sessionParams)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not create session: %v", err)
	}

	return id, nil
}

func IssueSession(ctx context.Context, dbQueries *dbgen.Queries, sessionID uuid.UUID, userID uuid.UUID, name string, email string, secretKey []byte) (string, error) {
	expiresAt, err := dbQueries.GetSessionExpiration(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("could not get session expiration: %v", err)
	}

	revalidateAt := time.Now().Add(1 * time.Minute).Unix()

	sessionClaims := SessionClaims{
		SessionID:    sessionID,
		UserID:       userID,
		Name:         name,
		Email:        email,
		RevalidateAt: revalidateAt,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionClaims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not create signed token: %v", err)
	}

	return signedToken, nil
}

func CheckPassword(password, hashStr, saltStr string) (bool, error) {
	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, fmt.Errorf("could not decode salt string: %v", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return false, fmt.Errorf("could not decode hash string: %v", err)
	}

	var (
		memory  uint32 = 64 * 1024
		iter    uint32 = 3
		threads uint8  = 2
		keyLen  uint32 = 32
	)

	test := argon2.IDKey([]byte(password), salt, iter, memory, threads, keyLen)

	return bytes.Equal(hash, test), nil
}
