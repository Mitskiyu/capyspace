package widget

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	userIdRaw := ctx.Value("user_id")
	userId, ok := userIdRaw.(string)
	if !ok {
		log.Printf("mismatched type for user_id: %T", userIdRaw)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Print(userId) // verify userId against spaceId

	spaceId, err := uuid.Parse(chi.URLParam(r, "spaceId"))
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	widget, err := h.service.createWidget(ctx, req.Id, spaceId, string(req.Type), req.XPos, req.YPos, req.Minimized, req.Data)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Print(widget)
}
