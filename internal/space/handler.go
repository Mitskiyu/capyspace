package space

import (
	"log"
	"net/http"
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
		log.Printf("space already exists for %v", userId)
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
