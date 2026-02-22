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
	UpdateWidgetData(ctx context.Context, arg sqlc.UpdateWidgetDataParams) (sqlc.Widget, error)
}

type service struct {
	store Store
}

func NewService(store Store) service {
	return service{
		store: store,
	}
}

func (s *service) createWidget(ctx context.Context, id, spaceId uuid.UUID, widgetType string, xPos, yPos int32, minimized bool, data json.RawMessage) (sqlc.Widget, error) {
	params := sqlc.CreateWidgetParams{
		ID:        id,
		SpaceID:   spaceId,
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
