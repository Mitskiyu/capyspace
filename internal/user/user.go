package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
)

type Store interface {
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	GetUserByUsername(ctx context.Context, username string) (sqlc.User, error)
}

type service struct {
	store Store
}

func NewService(store Store) service {
	return service{
		store: store,
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
