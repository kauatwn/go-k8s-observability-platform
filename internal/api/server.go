package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter cria o multiplexador e registra todas as rotas
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	home := NewHomeHandler()
	health := NewHealthHandler()

	// Rota raiz da API envolta pelo middleware de métricas
	mux.Handle("/", RecordMetricsMiddleware(http.HandlerFunc(home.Home)))

	// Rotas da API envoltas pelo middleware de métricas
	mux.Handle("/healthz", RecordMetricsMiddleware(http.HandlerFunc(health.Healthz)))
	mux.Handle("/ready", RecordMetricsMiddleware(http.HandlerFunc(health.Ready)))

	// Rota interna do Prometheus (usando o handler oficial da biblioteca)
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
