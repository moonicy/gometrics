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

func TestMetricsHandler_GetJSONMetricsByName(t *testing.T) {
	ctx := context.Background()
	defaultMemStorage := storage.NewMemStorage()
	presetMemStorage := func() Storage {
		mem := storage.NewMemStorage()
		mem.AddCounter(ctx, agent.Alloc, 22)
		mem.SetGauge(ctx, agent.Frees, 22)
		return mem
	}
	tests := []struct {
		name   string
		mem    Storage
		body   metrics.MetricName
		status int
	}{
		{name: "gauge not found", body: metrics.MetricName{ID: agent.Alloc, MType: agent.Gauge}, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "counter not found", body: metrics.MetricName{ID: agent.Frees, MType: agent.Counter}, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "response 200 for gauge", body: metrics.MetricName{ID: agent.Frees, MType: agent.Gauge}, mem: presetMemStorage(), status: http.StatusOK},
		{name: "response 200 for counter", body: metrics.MetricName{ID: agent.Alloc, MType: agent.Counter}, mem: presetMemStorage(), status: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(tt.body)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/value/", bytes.NewBuffer(out))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				storage: tt.mem,
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/value/", u.GetMetricValueByNameJSON)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func ExampleMetricsHandler_GetMetricValueByNameJSON() {
	// Инициализируем хранилище и добавляем метрику типа gauge.
	memStorage := storage.NewMemStorage()
	_ = memStorage.SetGauge(context.Background(), "Alloc", 12345.67)

	// Создаём новый MetricsHandler.
	mh := NewMetricsHandler(memStorage, nil, nil)

	// Подготавливаем тело запроса в формате JSON.
	metricName := metrics.MetricName{
		ID:    "Alloc",
		MType: metrics.Gauge,
	}
	requestBody, _ := json.Marshal(metricName)

	// Создаём новый HTTP-запрос с телом в формате JSON.
	req := httptest.NewRequest("POST", "/value/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаём Recorder для записи HTTP-ответа.
	rr := httptest.NewRecorder()

	// Вызываем обработчик.
	mh.GetMetricValueByNameJSON(rr, req)

	// Выводим статусный код и тело ответа.
	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Body: {"id":"Alloc","type":"gauge","value":12345.67}
}
