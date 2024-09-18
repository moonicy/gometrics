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

func TestUpdateMetrics_updateJSONMetrics(t *testing.T) {
	value := 11.1
	delta := int64(11)
	tests := []struct {
		status int
		name   string
		body   metrics.Metric
	}{
		{name: "response 200 for gauge", body: metrics.Metric{MetricName: metrics.MetricName{ID: agent.Alloc, MType: agent.Gauge}, Value: &value}, status: http.StatusOK},
		{name: "response 200 for counter", body: metrics.Metric{MetricName: metrics.MetricName{ID: agent.Frees, MType: agent.Counter}, Delta: &delta}, status: http.StatusOK},
		{name: "wrong type", body: metrics.Metric{MetricName: metrics.MetricName{ID: agent.Alloc, MType: "wrong"}, Delta: &delta}, status: http.StatusBadRequest},
		{name: "without name", body: metrics.Metric{MetricName: metrics.MetricName{ID: "", MType: agent.Gauge}, Value: &value}, status: http.StatusNotFound},
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
				storage: storage.NewMemStorage(),
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/", u.PostMetricUpdateJSON)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func ExampleMetricsHandler_PostMetricUpdateJSON() {
	// Инициализируем хранилище.
	memStorage := storage.NewMemStorage()

	// Создаём новый MetricsHandler.
	mh := NewMetricsHandler(memStorage, nil, nil)

	// Создаём метрику для обновления в формате JSON.
	metric := metrics.Metric{
		MetricName: metrics.MetricName{
			ID:    "Alloc",
			MType: metrics.Gauge,
		},
		Value: new(float64),
	}
	*metric.Value = 12345.67

	// Кодируем метрику в JSON.
	body, _ := json.Marshal(metric)

	// Создаём HTTP-запрос для обновления метрики.
	req := httptest.NewRequest("POST", "/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаём Recorder для записи ответа.
	rr := httptest.NewRecorder()

	// Вызываем обработчик.
	mh.PostMetricUpdateJSON(rr, req)

	// Выводим статусный код.
	fmt.Println("Status Code:", rr.Code)

	// Проверяем, что метрика была обновлена.
	value, err := memStorage.GetGauge(context.Background(), "Alloc")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Updated Gauge Alloc: %.2f\n", value)
	}

	// Output:
	// Alloc	gauge	12345.670000
	// Status Code: 200
	// Updated Gauge Alloc: 12345.67
}
