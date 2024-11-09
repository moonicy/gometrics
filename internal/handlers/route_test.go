package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockMetricsHandler struct{}

func (m *MockMetricsHandler) GetMetrics(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockMetricsHandler) GetMetricValueByNameJSON(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockMetricsHandler) GetMetricValueByName(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockMetricsHandler) PostMetricUpdateJSON(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockMetricsHandler) PostMetricUpdate(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockMetricsHandler) PostMetricsUpdatesJSON(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockMetricsHandler) GetPing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestNewRoute(t *testing.T) {
	mh := &MockMetricsHandler{}
	log := zap.NewExample().Sugar()
	defer log.Sync()

	cfg := config.ServerConfig{
		CryptoKey:     "test-crypto-key",
		HashKey:       "test-hash-key",
		TrustedSubnet: "192.168.1.0/24",
	}

	router := NewRoute(mh, log, cfg)

	tests := []struct {
		method     string
		target     string
		statusCode int
	}{
		{method: "GET", target: "/", statusCode: http.StatusOK},
		{method: "POST", target: "/value", statusCode: http.StatusOK},
		{method: "GET", target: "/value/counter/example", statusCode: http.StatusOK},
		{method: "POST", target: "/update", statusCode: http.StatusOK},
		{method: "POST", target: "/update/gauge/example/100", statusCode: http.StatusOK},
		{method: "POST", target: "/updates", statusCode: http.StatusForbidden},
		{method: "GET", target: "/ping", statusCode: http.StatusOK},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.target, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, tt.statusCode, resp.StatusCode, "Expected status code %d for %s %s, got %d", tt.statusCode, tt.method, tt.target, resp.StatusCode)
	}
}
