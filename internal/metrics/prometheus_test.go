package metrics

import "testing"

func TestMetrics(t *testing.T) {
	// Garante que os coletores globais foram inicializados na carga do pacote
	if httpRequestsTotal == nil || httpRequestDuration == nil {
		t.Fatal("expected Prometheus collectors to be initialized")
	}

	// Simula o registro de uma requisição válida
	RecordHTTPRequest("/{$}", "GET", 200, 0.025)

	// Protege o teste contra panic caso Register() seja chamado mais de uma vez
	// no mesmo processo de testes (comportamento padrão do registry global do Prometheus).
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Register recovered: %v", r)
		}
	}()

	Register()
}
