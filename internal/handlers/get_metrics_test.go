package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/storage"
)

func TestMetricsHandler_GetMetrics(t *testing.T) {
	ctx := context.Background()
	const bodyWait = "Alloc: 22\nFrees: 22\n"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	mem := storage.NewMemStorage()
	mem.AddCounter(ctx, agent.Alloc, 22)
	mem.SetGauge(ctx, agent.Frees, 22)

	u := &MetricsHandler{
		storage: mem,
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

func ExampleMetricsHandler_GetMetrics() {
	// Инициализируем хранилище и добавляем метрики.
	memStorage := storage.NewMemStorage()
	_ = memStorage.SetGauge(context.Background(), "Alloc", 12345.67)
	_ = memStorage.AddCounter(context.Background(), "PollCount", 42)

	// Создаём новый MetricsHandler.
	mh := NewMetricsHandler(memStorage, nil, nil)

	// Создаём новый HTTP-запрос.
	req := httptest.NewRequest("GET", "/", nil)

	// Создаём Recorder для записи HTTP-ответа.
	rr := httptest.NewRecorder()

	// Вызываем обработчик.
	mh.GetMetrics(rr, req)

	// Выводим статусный код и тело ответа.
	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Body:")
	fmt.Println(strings.TrimSpace(rr.Body.String()))

	// Output:
	// Status Code: 200
	// Body:
	// PollCount: 42
	// Alloc: 12345.67
}
