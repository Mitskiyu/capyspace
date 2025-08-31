package auth

import (
	"context"
	"fmt"

	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/google/uuid"
)

type Store interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
}

type service struct {
	store Store
}

func NewService(store Store) service {
	return service{
		store: store,
	}
}

func (s *service) register(ctx context.Context, email, password string) (*sqlc.User, error) {
	password = hashPassword(password)
	params := sqlc.CreateUserParams{
		ID:       uuid.New(),
		Email:    email,
		Password: password,
	}

	user, err := s.store.CreateUser(ctx, params)
	if err != nil {
		return &sqlc.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}
