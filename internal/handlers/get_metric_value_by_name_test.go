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

func TestMetricsHandler_GetMetricsByName(t *testing.T) {
	ctx := context.Background()
	defaultMemStorage := storage.NewMemStorage()
	presetMemStorage := func() Storage {
		mem := storage.NewMemStorage()
		err := mem.AddCounter(ctx, agent.Alloc, 22)
		if err != nil {
			log.Fatal(err)
		}
		err = mem.SetGauge(ctx, agent.Frees, 22)
		if err != nil {
			log.Fatal(err)
		}
		return mem
	}
	tests := []struct {
		name    string
		tpMet   string
		nameMet string
		mem     Storage
		status  int
	}{
		{name: "gauge not found", tpMet: agent.Gauge, nameMet: agent.Alloc, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "counter not found", tpMet: agent.Counter, nameMet: agent.Frees, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "response 200 for gauge", tpMet: agent.Gauge, nameMet: agent.Frees, mem: presetMemStorage(), status: http.StatusOK},
		{name: "response 200 for counter", tpMet: agent.Counter, nameMet: agent.Alloc, mem: presetMemStorage(), status: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/value/%s/%s", tt.tpMet, tt.nameMet)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				storage: tt.mem,
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/value/{type}/{name}", u.GetMetricValueByName)
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

func ExampleMetricsHandler_GetMetricValueByName() {
	// Инициализируем хранилище и добавляем метрику типа gauge.
	memStorage := storage.NewMemStorage()
	_ = memStorage.SetGauge(context.Background(), "Alloc", 12345.67)

	// Создаём новый MetricsHandler.
	mh := NewMetricsHandler(memStorage, nil, nil)

	// Создаём новый HTTP-запрос.
	req := httptest.NewRequest("GET", "/value/gauge/Alloc", nil)

	// Устанавливаем параметры URL в контексте chi.
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("name", "Alloc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Создаём Recorder для записи HTTP-ответа.
	rr := httptest.NewRecorder()

	// Вызываем обработчик.
	mh.GetMetricValueByName(rr, req)

	// Выводим статусный код и тело ответа.
	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Body: 12345.67
}
