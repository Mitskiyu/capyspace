package widget

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/google/uuid"
)

type Store interface {
	CreateWidget(ctx context.Context, arg sqlc.CreateWidgetParams) (sqlc.Widget, error)
	UpdateWidget(ctx context.Context, arg sqlc.UpdateWidgetParams) (sqlc.Widget, error)
}

type service struct {
	store Store
}

func NewService(store Store) service {
	return service{
		store: store,
	}
}

func (s *service) createWidget(ctx context.Context, id, spaceID uuid.UUID, widgetType string, xPos, yPos int32, minimized bool, data json.RawMessage) (sqlc.Widget, error) {
	params := sqlc.CreateWidgetParams{
		ID:        id,
		SpaceID:   spaceID,
		Type:      widgetType,
		XPos:      xPos,
		YPos:      yPos,
		Minimized: minimized,
		Data:      data,
	}

	widget, err := s.store.CreateWidget(ctx, params)
	if err != nil {
		return sqlc.Widget{}, fmt.Errorf("failed to create widget: %w", err)
	}

	return widget, nil
}

func (s *service) updateWidget(ctx context.Context, id uuid.UUID, xPos, yPos int32, minimized bool, data json.RawMessage) (sqlc.Widget, error) {
	params := sqlc.UpdateWidgetParams{
		ID:        id,
		XPos:      xPos,
		YPos:      yPos,
		Minimized: minimized,
		Data:      data,
	}

	widget, err := s.store.UpdateWidget(ctx, params)
	if err != nil {
		return sqlc.Widget{}, fmt.Errorf("failed to update widget: %w", err)
	}

	return widget, nil
}
