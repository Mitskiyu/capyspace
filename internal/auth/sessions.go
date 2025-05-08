package auth

import (
	"context"
	"fmt"
	"time"

	dbgen "github.com/Mitskiyu/capyspace/internal/database/sqlc"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Creates a 90-day database session in the db and returns the session's id
func createSession(ctx context.Context, dbq *dbgen.Queries, userID uuid.UUID) (uuid.UUID, error) {
	sessionParams := dbgen.CreateSessionParams{
		ID:        uuid.New(),
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	id, err := dbq.CreateSession(ctx, sessionParams)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not create session: %w", err)
	}

	return id, nil
}

// Issues a JWT session and returns the signed token as a string
// Creates and signs a JWT token for an existing session in the db
func issueSession(ctx context.Context, dbq *dbgen.Queries, sessionID uuid.UUID, userID uuid.UUID, name string, email string, sk []byte) (string, error) {
	// Gets the database session's expiration time
	expiresAt, err := dbq.GetSessionExpiration(ctx, sessionID)
	if err != nil {
		return "", err
	}

	// Set revalidateAt to a minute from now, so the token has to be revalidated every minute
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

	// Signs the JWT with the server's secret key, so the token can be validated by the server
	signedToken, err := token.SignedString(sk)
	if err != nil {
		return "", fmt.Errorf("could not create signed token: %w", err)
	}

	return signedToken, nil
}

// Revalidates a JWT session and returns the new token and its claims
func RevalidateSession(ctx context.Context, dbq *dbgen.Queries, sk []byte, tokenStr string) (string, *SessionClaims, error) {
	// Parses the token and validates it using the secret key
	token, err := jwt.ParseWithClaims(tokenStr, &SessionClaims{}, func(token *jwt.Token) (any, error) {
		return sk, nil
	})
	if err != nil || !token.Valid {
		return "", nil, fmt.Errorf("token is invalid")
	}

	// Checks the token's claims are of the expected type, making sure they are valid
	claims, ok := token.Claims.(*SessionClaims)
	if !ok {
		return "", nil, fmt.Errorf("invalid claims")
	}

	// Checks if the token needs to be revalidated
	if time.Now().Unix() < claims.RevalidateAt {
		return tokenStr, claims, nil
	}

	session, err := dbq.GetSession(ctx, claims.SessionID)
	if err != nil {
		return "", nil, err
	}

	// Don't revalidate the JWT session if the database session is revoked or expired
	if session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return "", nil, fmt.Errorf("session is expired or revoked")
	}

	user, err := dbq.GetUser(ctx, claims.UserID)
	if err != nil {
		return "", nil, err
	}

	// Creates new session claims for the refreshed JWT session
	newClaims := SessionClaims{
		SessionID:    session.ID,
		UserID:       user.ID,
		Name:         user.Name.String,
		Email:        user.Email,
		RevalidateAt: time.Now().Add(1 * time.Minute).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(session.ExpiresAt),
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signedToken, err := newToken.SignedString(sk)
	if err != nil {
		return "", nil, err
	}

	return signedToken, &newClaims, nil
}
