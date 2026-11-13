package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	handler := NewHomeHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.Home(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	expected := "Welcome to Go K8s Observability Platform!"
	if resp["message"] != expected {
		t.Errorf("expected message '%s', got '%s'", expected, resp["message"])
	}
}
