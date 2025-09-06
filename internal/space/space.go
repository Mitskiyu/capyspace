package space

import (
	"context"
	"fmt"

	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Store interface {
	CreateSpace(ctx context.Context, arg sqlc.CreateSpaceParams) (sqlc.Space, error)
}

type service struct {
	store Store
}

func NewService(store Store) service {
	return service{
		store: store,
	}
}

func (s *service) createSpace(ctx context.Context, userId string) (bool, sqlc.Space, error) {
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		return false, sqlc.Space{}, fmt.Errorf("failed to parse uuid: %w", err)
	}

	params := sqlc.CreateSpaceParams{
		ID:     uuid.New(),
		UserID: parsedId,
	}

	space, err := s.store.CreateSpace(ctx, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return false, sqlc.Space{}, nil
		}
		return false, sqlc.Space{}, fmt.Errorf("failed to create space: %w", err)
	}

	return true, space, nil
}
