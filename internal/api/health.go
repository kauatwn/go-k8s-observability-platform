package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

type HealthHandler struct {
	isReady atomic.Bool
}

func NewHealthHandler() *HealthHandler {
	h := &HealthHandler{}

	go func() {
		// Simula um tempo de inicialização do serviço (3 segundos de aquecimento)
		time.Sleep(3 * time.Second)
		h.isReady.Store(true)
		slog.Info("Service is now ready to accept traffic")
	}()

	return h
}

// Healthz retorna o status de saúde do processo (Liveness Probe)
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

// Ready retorna o status de prontidão do serviço (Readiness Probe)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Se o serviço ainda não estiver pronto, retorna 503 Service Unavailable
	if !h.isReady.Load() {
		respondJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "NOT READY"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "READY"})
}

// respondJSON centraliza a escrita dos cabeçalhos, status code e o tratamento de erro do encoder JSON
func respondJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
