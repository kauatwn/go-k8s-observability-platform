package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HomeHandler agrupa as dependências das rotas
type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to Go K8s Observability Platform!",
	}); err != nil {
		slog.Error("Failed to encode home response", "error", err)
	}
}
