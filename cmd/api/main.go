package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-k8s-observability-platform/internal/api"
	"go-k8s-observability-platform/internal/metrics"
)

func main() {
	// Configura o logger estruturado em formato JSON enviando para a saída padrão (stdout)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting application", slog.String("env", "development"))

	// Registra as métricas da aplicação no coletor do Prometheus
	metrics.Register()

	// Configura as rotas da API
	mux := api.SetupRouter()

	// Configura o servidor HTTP com timeouts estritos de segurança
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,  // Mitiga ataques de Slowloris
		ReadTimeout:       5 * time.Second,  // Limite para leitura do corpo da requisição
		WriteTimeout:      10 * time.Second, // Limite para escrita da resposta
	}

	// Configura o contexto para capturar sinais de interrupção (Ctrl+C) e SIGTERM do Kubernetes
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Inicia o servidor HTTP em uma goroutine separada para não bloquear a execução principal
	go func() {
		slog.Info("Server is listening", slog.String("addr", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Aguarda o contexto ser cancelado (sinal de interrupção ou SIGTERM)
	<-ctx.Done()
	slog.Info("Received shutdown signal, starting graceful shutdown...")

	// Cria um contexto com tempo limite de 10s para drenar as requisições em andamento
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Encerra o servidor graciosamente
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown due to error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped gracefully. Bye!")
}
