package server

import "net/http"

func (s Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
