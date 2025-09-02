package auth

import (
	"log"
	"net/http"

	"github.com/Mitskiyu/capyspace/internal/user"
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

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	req, _, err := util.Decode[RegisterReq](r)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if _, err := h.service.register(ctx, req.Email, req.Password); err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	req, _, err := util.Decode[LoginReq](r)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, sessionId, err := h.service.login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Could not log in", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		Path:     "/",
	}

	res := LoginRes{
		Message: "Login successful",
		User: user.Info{
			Id:    u.ID.String(),
			Email: u.Email,
		},
	}

	http.SetCookie(w, cookie)
	if err := util.Encode(w, http.StatusOK, res); err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
	}
}
