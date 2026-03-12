package space

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/util"
	"github.com/Mitskiyu/capyspace/internal/widget"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateSpace(w http.ResponseWriter, r *http.Request) {
	userIdRaw := r.Context().Value("user_id")
	userID, ok := userIdRaw.(string)
	if !ok {
		log.Printf("mismatched type for user_id: %T", userIdRaw)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	created, _, err := h.service.createSpace(r.Context(), userID)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !created {
		log.Println("space already exists for: &s", userID)
		http.Error(w, "Space already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetSpace(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	ctx := r.Context()
	found, space, widgets, err := h.service.getSpace(ctx, username)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "Space not found", http.StatusNotFound)
		return
	}

	widgetRes := make([]widget.Widget, len(widgets))
	for i, wgt := range widgets {
		widgetRes[i] = widget.Widget{
			ID:        wgt.ID,
			Type:      widget.Type(wgt.Type),
			XPos:      wgt.XPos,
			YPos:      wgt.YPos,
			Minimized: wgt.Minimized,
			Data:      wgt.Data,
		}
	}

	res := SpaceRes{
		ID:        space.ID.String(),
		IsPrivate: space.IsPrivate,
		Widgets:   widgetRes,
	}

	if space.IsPrivate {
		if space.UserID.String() == ctx.Value("user_id") {
			fmt.Print(ctx.Value("user_id"))
			util.Encode(w, http.StatusOK, res)
			return
		} else {
			http.Error(w, "Space is private", http.StatusForbidden)
			return
		}
	}

	util.Encode(w, http.StatusOK, res)
}
