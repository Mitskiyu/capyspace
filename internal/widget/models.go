package widget

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Type string

const (
	TypeStickyNote Type = "sticky-note"
)

type CreateWidgetReq struct {
	Id        uuid.UUID       `json:"id"`
	Type      Type            `json:"type"`
	XPos      int32           `json:"x_pos"`
	YPos      int32           `json:"y_pos"`
	Minimized bool            `json:"minimized"`
	Data      json.RawMessage `json:"data"`
}

func (c CreateWidgetReq) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	switch c.Type {
	case TypeStickyNote:
	default:
		problems["type"] = "invalid widget type"
	}

	return problems
}
