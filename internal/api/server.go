// Package api implementa os handlers HTTP, o roteamento e os middlewares
// de observabilidade da aplicação.
package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter cria o multiplexador HTTP, registra os endpoints da aplicação
// e envolve o roteador inteiro com o middleware de métricas e logs.
// Retorna a interface http.Handler para manter compatibilidade com http.Server.
func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	home := NewHomeHandler()
	health := NewHealthHandler()

	// O sufixo {$} exige correspondência exata com a raiz "/",
	// garantindo que URLs desconhecidas retornem 404 Not Found em vez de cair aqui.
	mux.HandleFunc("/{$}", home.Home)

	// Endpoints utilizados pelas probes de ciclo de vida do Kubernetes
	mux.HandleFunc("/healthz", health.Healthz) // Liveness probe
	mux.HandleFunc("/ready", health.Ready)     // Readiness probe

	// Endpoint padrão onde o Prometheus realiza o scrape das métricas
	mux.Handle("/metrics", promhttp.Handler())

	return RecordMetricsMiddleware(mux)
}
