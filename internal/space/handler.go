package space

import (
	"log"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/util"
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
	userId, ok := userIdRaw.(string)
	if !ok {
		log.Printf("mismatched type for user_id: %T", userIdRaw)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	created, _, err := h.service.createSpace(r.Context(), userId)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !created {
		log.Println("space already exists for: &s", userId)
		http.Error(w, "Space already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetSpace(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	found, space, err := h.service.getSpace(r.Context(), username)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "Space not found", http.StatusNotFound)
		return
	}

	if space.IsPrivate {
		http.Error(w, "Space is private", http.StatusForbidden)
		return
	}

	res := SpaceRes{
		Id:        space.ID.String(),
		IsPrivate: space.IsPrivate,
	}

	util.Encode(w, http.StatusOK, res)
}
