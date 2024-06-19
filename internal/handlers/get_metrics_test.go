package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_GetMetrics(t *testing.T) {
	const bodyWait = "alloc: 22\nfrees: 22\n"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mem := storage.NewMemStorage()
	mem.AddCounter(agent.Alloc, 22)
	mem.SetGauge(agent.Frees, 22)

	u := &MetricsHandler{
		mem: mem,
	}

	rec := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/", u.GetMetrics)
	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}
	bodyString := string(bodyBytes)
	if bodyString != bodyWait {
		t.Errorf("expected: %s\ngot: \n%s", bodyWait, bodyString)
	}

}
