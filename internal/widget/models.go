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

type Widget struct {
	ID        uuid.UUID       `json:"id"`
	Type      Type            `json:"type"`
	XPos      int32           `json:"x_pos"`
	YPos      int32           `json:"y_pos"`
	Minimized bool            `json:"minimized"`
	Data      json.RawMessage `json:"data"`
}

type StickyNoteData struct {
	Text string `json:"text"`
}

type CreateWidgetReq struct {
	Widget
}

type UpdateWidgetReq struct {
	Widget
	ID   uuid.UUID `json:"-"` // getting id from url, override
	Type Type      `json:"-"` // not updating type
}

func (w Widget) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	switch w.Type {
	case TypeStickyNote:
		var data StickyNoteData
		if err := json.Unmarshal(w.Data, &data); err != nil {
			problems["data"] = "invalid sticky note data"
		}
	default:
		problems["type"] = "invalid widget type"
	}

	return problems
}
