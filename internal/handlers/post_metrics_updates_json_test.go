package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_UpdatesJSONMetrics(t *testing.T) {
	value := 11.1
	delta := int64(11)
	tests := []struct {
		status int
		name   string
		body   []metrics.Metric
	}{
		{name: "response 200",
			body: []metrics.Metric{
				{MetricName: metrics.MetricName{ID: agent.Alloc, MType: agent.Gauge}, Value: &value},
				{MetricName: metrics.MetricName{ID: agent.Frees, MType: agent.Gauge}, Value: &value},
				{MetricName: metrics.MetricName{ID: agent.PollCount, MType: agent.Counter}, Delta: &delta},
			},
			status: http.StatusOK},
		{name: "wrong type",
			body:   []metrics.Metric{{MetricName: metrics.MetricName{ID: agent.Alloc, MType: "wrong"}, Delta: &delta}},
			status: http.StatusBadRequest},
		{name: "without name",
			body:   []metrics.Metric{{MetricName: metrics.MetricName{ID: "", MType: agent.Gauge}, Value: &value}},
			status: http.StatusNotFound},
		{name: "with empty body",
			body:   []metrics.Metric{},
			status: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(tt.body)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/updates/", bytes.NewBuffer(out))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				storage: storage.NewMemStorage(),
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/updates/", u.PostMetricsUpdatesJSON)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}
