package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Store interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	GetUserByUsername(ctx context.Context, username string) (sqlc.User, error)
}

type Cache interface {
	SetSession(ctx context.Context, sessionId, userId string, exp time.Duration) error
	GetSession(ctx context.Context, sessionId string) (string, time.Duration, error)
	UpdateSessionTTL(ctx context.Context, sessionId string, exp time.Duration) error
}

type service struct {
	store Store
	cache Cache
}

func NewService(store Store, cache Cache) service {
	return service{
		store: store,
		cache: cache,
	}
}

func (s *service) checkEmail(ctx context.Context, email string) (bool, error) {
	_, err := s.store.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, err
	}
}

func (s *service) checkUsername(ctx context.Context, username string) (bool, error) {
	_, err := s.store.GetUserByUsername(ctx, username)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, err
	}
}

func (s *service) register(ctx context.Context, email, password, username string) (bool, *sqlc.User, error) {
	password = hashPassword(password)
	params := sqlc.CreateUserParams{
		ID:          uuid.New(),
		Email:       email,
		Password:    password,
		Username:    username,
		DisplayName: username,
	}

	user, err := s.store.CreateUser(ctx, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("failed to create user: %w", err)
	}

	return true, &user, nil
}

func (s *service) login(ctx context.Context, email, password string) (bool, *sqlc.User, string, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil, "", nil
	case err != nil:
		return false, nil, "", fmt.Errorf("failed to find user: %w", err)
	}

	equal, err := comparePassword(user.Password, password)
	if err != nil {
		return false, &user, "", fmt.Errorf("failed to compare password: %w", err)
	}

	if !equal {
		return false, &user, "", nil
	}

	sessionId := createSessionId()
	exp := 24 * 30 * time.Hour
	if err := s.cache.SetSession(ctx, sessionId, user.ID.String(), exp); err != nil {
		return false, &sqlc.User{}, "", fmt.Errorf("failed to set session for user %s: %v", user.ID.String(), err)
	}

	return true, &user, sessionId, nil
}

func (s *service) sessionMiddleware(ctx context.Context, sessionId string) (string, bool, error) {
	if err := validSessionId(sessionId); err != nil {
		return "", false, fmt.Errorf("failed to validate session id %s: %w", sessionId, err)
	}

	userId, ttl, err := s.cache.GetSession(ctx, sessionId)
	if err != nil {
		return "", false, fmt.Errorf("failed to get session for %s: %w", sessionId, err)
	}

	expiring := ttl < 7*24*time.Hour
	if expiring {
		if err := s.cache.UpdateSessionTTL(ctx, sessionId, 30*24*time.Hour); err != nil {
			return "", expiring, fmt.Errorf("failed to update session ttl for %s: %w", sessionId, err)
		}

		return userId, expiring, nil
	}

	return userId, false, nil
}
