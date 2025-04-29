package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/google/uuid"
)

func GenerateVerificationToken() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("could not generate token: %v", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func CreateVerificationToken(ctx context.Context, dbQueries *dbgen.Queries, email string) (string, error) {
	token, err := GenerateVerificationToken()
	if err != nil {
		return "", err
	}

	tokenParams := dbgen.CreateVerificationTokenParams{
		ID:        uuid.New(),
		Email:     email,
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := dbQueries.CreateVerificationToken(ctx, tokenParams); err != nil {
		return "", fmt.Errorf("could not create token: %v", err)
	}
	return token, nil
}

func CheckVerificationToken(ctx context.Context, dbQueries *dbgen.Queries, email string, token string) (bool, error) {
	tokenParams := dbgen.GetValidVerificationTokenParams{
		Email: email,
		Token: token,
	}

	vt, err := dbQueries.GetValidVerificationToken(ctx, tokenParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("could not get valid token: %v", err)
	}

	// set the token's <used> to true
	if err := dbQueries.SetUsedVerificationToken(ctx, vt.ID); err != nil {
		return false, fmt.Errorf("could not set used token: %v", err)
	}

	return true, nil
}
