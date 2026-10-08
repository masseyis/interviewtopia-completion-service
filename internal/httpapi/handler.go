package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/masseyis/interviewtopia-completion-service/pkg/evidence"
)

// Dependencies are supplied by the executable and can be replaced directly in
// tests. The starter health route does not use them; they are ready for the
// exercise implementation.
type Dependencies struct {
	RegistryURL string
	TrustStore  evidence.TrustStore
	HTTPClient  *http.Client
}

type handler struct {
	dependencies Dependencies
	mux          *http.ServeMux
}

func NewHandler(dependencies Dependencies) http.Handler {
	if dependencies.HTTPClient == nil {
		dependencies.HTTPClient = http.DefaultClient
	}

	h := &handler{
		dependencies: dependencies,
		mux:          http.NewServeMux(),
	}
	h.mux.HandleFunc("/healthz", h.health)
	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
