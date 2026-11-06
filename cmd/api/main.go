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

	// Executa o ciclo de vida da aplicação e trata qualquer erro fatal na raiz
	if err := run(); err != nil {
		slog.Error("Application failed to run", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("Starting application", slog.String("env", "development"))

	// Registra as métricas da aplicação no coletor do Prometheus
	metrics.Register()

	// Configura as rotas da API
	mux := api.SetupRouter()

	// Configura o servidor HTTP com timeouts estritos de segurança
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,   // Mitiga ataques de Slowloris
		ReadTimeout:       5 * time.Second,   // Limite para leitura do corpo da requisição
		WriteTimeout:      10 * time.Second,  // Limite para escrita da resposta
		IdleTimeout:       120 * time.Second, // Tempo máximo de inatividade antes de fechar a conexão (keep-alive)
	}

	// Configura o contexto para capturar sinais de interrupção (Ctrl+C) e SIGTERM do Kubernetes
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Canal buffered para capturar erros de inicialização ou execução do servidor HTTP
	serverErrors := make(chan error, 1)

	// Inicia o servidor HTTP em uma goroutine separada para não bloquear a execução principal
	go func() {
		slog.Info("Server is listening", slog.String("addr", server.Addr))
		serverErrors <- server.ListenAndServe()
	}()

	// Aguarda falha inesperada no servidor OU sinal de término do sistema operacional/Kubernetes
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil

	case <-ctx.Done():
		slog.Info("Received shutdown signal, starting graceful shutdown...")
	}

	// Cria um contexto com tempo limite de 10s para drenar as requisições em andamento
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tenta encerrar o servidor graciosamente
	err := server.Shutdown(shutdownCtx)
	if err == nil {
		slog.Info("Server stopped gracefully. Bye!")
		return nil
	}

	// Força o fechamento das conexões ativas caso o graceful shutdown falhe
	if closeErr := server.Close(); closeErr != nil {
		return errors.Join(err, closeErr)
	}

	return err
}
