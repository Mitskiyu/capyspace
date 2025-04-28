package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/google/uuid"
)

func GenerateToken() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("could not generate token: %v", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func CreateToken(ctx context.Context, dbQueries *dbgen.Queries, email string) (string, error) {
	token, err := GenerateToken()
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
		return "", fmt.Errorf("database error: %v", err)
	}
	return token, nil
}
