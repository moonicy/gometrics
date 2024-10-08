package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/storage"
)

func TestUpdateMetrics_updateMetrics(t *testing.T) {
	tests := []struct {
		name    string
		tpMet   string
		nameMet string
		valMet  string
		status  int
	}{
		{name: "response 200 for gauge", tpMet: agent.Gauge, nameMet: agent.Alloc, valMet: "11.1", status: http.StatusOK},
		{name: "response 200 for counter", tpMet: agent.Counter, nameMet: agent.Frees, valMet: "11", status: http.StatusOK},
		{name: "wrong type", tpMet: "wrong", nameMet: agent.Alloc, valMet: "11", status: http.StatusBadRequest},
		{name: "without name", tpMet: agent.Gauge, nameMet: "", valMet: "11", status: http.StatusNotFound},
		{name: "value for gauge not float", tpMet: agent.Gauge, nameMet: agent.Frees, valMet: "str", status: http.StatusBadRequest},
		{name: "value for counter not int", tpMet: agent.Counter, nameMet: agent.Alloc, valMet: "11.1", status: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", tt.tpMet, tt.nameMet, tt.valMet)
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				storage: storage.NewMemStorage(),
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", u.PostMetricUpdate)
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

func ExampleMetricsHandler_PostMetricUpdate() {
	// Инициализируем хранилище.
	memStorage := storage.NewMemStorage()

	// Создаём новый MetricsHandler.
	mh := NewMetricsHandler(memStorage, nil, nil)

	// Создаём HTTP-запрос для обновления метрики типа gauge.
	req := httptest.NewRequest("POST", "/update/gauge/Alloc/12345.67", nil)

	// Устанавливаем параметры URL в контексте chi.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("name", "Alloc")
	rctx.URLParams.Add("value", "12345.67")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Создаём Recorder для записи ответа.
	rr := httptest.NewRecorder()

	// Вызываем обработчик.
	mh.PostMetricUpdate(rr, req)

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
	// Alloc	12345.67	gauge
	// Status Code: 200
	// Updated Gauge Alloc: 12345.67
}
