package user

import (
	"log"
	"net/http"

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

func (h *handler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	req, _, err := util.Decode[EmailReq](r)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exists, err := h.service.checkEmail(ctx, req.Email)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := util.Encode(w, http.StatusOK, Exists{Exists: exists}); err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
	}
}

func (h *handler) CheckUsername(w http.ResponseWriter, r *http.Request) {
	req, _, err := util.Decode[UsernameReq](r)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exists, err := h.service.checkUsername(ctx, req.Username)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := util.Encode(w, http.StatusOK, Exists{Exists: exists}); err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
	}
}
