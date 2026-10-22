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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	}); err != nil {
		slog.Error("Failed to encode healthz response", "error", err)
	}
}

// Ready retorna o status de prontidão do serviço (Readiness Probe)
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !h.isReady.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": "NOT READY",
		}); err != nil {
			slog.Error("Failed to encode not ready response", "error", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "READY",
	}); err != nil {
		slog.Error("Failed to encode ready response", "error", err)
	}
}
