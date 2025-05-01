package auth

import (
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
	_, err := dbQueries.GetUsedVerficationCode(ctx, email)
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

	userParams := dbgen.CreateUserParams{
		ID:            uuid.New(),
		Name:          sql.NullString{String: "", Valid: false},
		Email:         email,
		Password:      sql.NullString{String: hashStr, Valid: true},
		EmailVerified: time.Now(),
	}

	id, err := dbQueries.CreateUser(ctx, userParams)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not create user: %v", err)
	}

	return id, nil
}
