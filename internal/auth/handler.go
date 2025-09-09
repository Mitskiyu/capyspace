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

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	req, _, err := util.Decode[RegisterReq](r)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if created, _, err := h.service.register(ctx, req.Email, req.Password, req.Username); err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !created {
		log.Printf("Attempted to register existing email: %s", req.Email)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
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
	ok, u, sessionId, err := h.service.login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("%v at %s", err, r.URL.Path)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !ok {
		res := LoginRes{
			Message: "Login failed",
			User: user.Info{
				Id:    "Invalid",
				Email: req.Email,
			},
		}
		if err := util.Encode(w, http.StatusUnauthorized, res); err != nil {
			log.Printf("%v at %s", err, r.URL.Path)
		}
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
