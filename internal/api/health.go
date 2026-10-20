package api

import (
	"encoding/json"
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
		// Simula um tempo de inicialização do serviço (3 segundos)
		time.Sleep(3 * time.Second)
		h.isReady.Store(true)
	}()
	return h
}

// Healthz retorna o status de saúde do serviço
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	})
	if err != nil {
		return
	}
}

// Ready retorna o status de prontidão do serviço
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !h.isReady.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		err := json.NewEncoder(w).Encode(map[string]string{
			"status": "NOT READY",
		})
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "READY",
	})
	if err != nil {
		return
	}
}
