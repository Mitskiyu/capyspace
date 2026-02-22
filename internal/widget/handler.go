package widget

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/database/sqlc"
	"github.com/Mitskiyu/capyspace/internal/util"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateWidget(w http.ResponseWriter, r *http.Request) {
	req, _, err := util.Decode[CreateWidgetReq](r)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	spaceRaw := ctx.Value("space")
	space, ok := spaceRaw.(*sqlc.Space)
	if !ok {
		log.Printf("mismatched type for space: %T", spaceRaw)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	widget, err := h.service.createWidget(ctx, req.Id, space.ID, string(req.Type), req.XPos, req.YPos, req.Minimized, req.Data)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Print(widget)
}
