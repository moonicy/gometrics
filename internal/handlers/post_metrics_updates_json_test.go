package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
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
			defer func() {
				if err = resp.Body.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func ExampleMetricsHandler_PostMetricsUpdatesJSON() {
	// Инициализируем хранилище.
	memStorage := storage.NewMemStorage()

	// Создаём новый MetricsHandler.
	mh := NewMetricsHandler(memStorage, nil, nil)

	// Создаём несколько метрик для обновления в формате JSON.
	metricsToUpdate := []metrics.Metric{
		{
			MetricName: metrics.MetricName{
				ID:    "Alloc",
				MType: metrics.Gauge,
			},
			Value: new(float64),
		},
		{
			MetricName: metrics.MetricName{
				ID:    "PollCount",
				MType: metrics.Counter,
			},
			Delta: new(int64),
		},
	}

	*metricsToUpdate[0].Value = 12345.67
	*metricsToUpdate[1].Delta = 42

	// Кодируем метрики в JSON.
	body, _ := json.Marshal(metricsToUpdate)

	// Создаём HTTP-запрос для обновления метрик.
	req := httptest.NewRequest("POST", "/updates", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаём Recorder для записи ответа.
	rr := httptest.NewRecorder()

	// Вызываем обработчик.
	mh.PostMetricsUpdatesJSON(rr, req)

	// Выводим статусный код.
	fmt.Println("Status Code:", rr.Code)

	// Проверяем, что метрики были обновлены.
	valueGauge, err := memStorage.GetGauge(context.Background(), "Alloc")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Updated Gauge Alloc: %.2f\n", valueGauge)
	}

	valueCounter, err := memStorage.GetCounter(context.Background(), "PollCount")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Updated Counter PollCount: %d\n", valueCounter)
	}

	// Output:
	// Status Code: 200
	// Updated Gauge Alloc: 12345.67
	// Updated Counter PollCount: 42
}
