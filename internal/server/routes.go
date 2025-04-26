package server

import (
	"net/http"
)

func (s Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	successResponse(w, http.StatusOK, "OK")
}
