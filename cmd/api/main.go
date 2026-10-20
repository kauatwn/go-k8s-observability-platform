package main

import (
	"go-k8s-observability-platform/internal/api"
	"go-k8s-observability-platform/internal/metrics"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting application",
		slog.String("env", "development"),
	)

	// Registra as métricas no Prometheus
	metrics.Register()

	// Configura as rotas da API
	mux := api.SetupRouter()

	// Configura o servidor HTTP com timeouts de segurança
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Server running",
		slog.String("addr", ":8080"),
	)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
