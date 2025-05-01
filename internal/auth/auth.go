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
