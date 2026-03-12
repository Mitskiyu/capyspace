package space

import "github.com/Mitskiyu/capyspace/internal/widget"

type SpaceRes struct {
	ID        string          `json:"id"`
	IsPrivate bool            `json:"is_private"`
	Widgets   []widget.Widget `json:"widgets"`
}
