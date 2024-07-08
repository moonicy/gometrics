package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/agent"
	l2p "github.com/moonicy/gometrics/internal/literaltopointer"
	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetrics_updateJSONMetrics(t *testing.T) {
	tests := []struct {
		status int
		name   string
		body   metrics.Metrics
	}{
		{name: "response 200 for gauge", body: metrics.Metrics{MetricName: metrics.MetricName{ID: agent.Alloc, MType: agent.Gauge}, Value: l2p.NewFloat(11.1)}, status: http.StatusOK},
		{name: "response 200 for counter", body: metrics.Metrics{MetricName: metrics.MetricName{ID: agent.Frees, MType: agent.Counter}, Delta: l2p.NewInt(11)}, status: http.StatusOK},
		{name: "wrong type", body: metrics.Metrics{MetricName: metrics.MetricName{ID: agent.Alloc, MType: "wrong"}, Delta: l2p.NewInt(11)}, status: http.StatusBadRequest},
		{name: "without name", body: metrics.Metrics{MetricName: metrics.MetricName{ID: "", MType: agent.Gauge}, Value: l2p.NewFloat(11)}, status: http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(tt.body)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/update/", bytes.NewBuffer(out))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				mem: storage.NewMemStorage(),
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/", u.UpdateJSONMetrics)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}
