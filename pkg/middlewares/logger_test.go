package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithLogging(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, test logging!"))
	})

	loggingMiddleware := WithLogging(*logger)(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	startTime := time.Now()

	loggingMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, test logging!", rec.Body.String())

	assert.Equal(t, 2, logs.Len())

	firstLog := logs.All()[0]

	assert.Contains(t, firstLog.Message, "uri /test method GET duration")
	assert.WithinDuration(t, startTime, startTime.Add(time.Since(startTime)), time.Millisecond)

	secondLog := logs.All()[1]
	assert.Equal(t, "status 200 responseSize 20", secondLog.Message)
}

func Test_loggingResponseWriter_Write(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Hello, world!"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	lrw := &loggingResponseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	handler.ServeHTTP(lrw, req)

	assert.Equal(t, http.StatusCreated, lrw.statusCode)
	assert.Equal(t, len("Hello, world!"), lrw.responseSize)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Hello, world!", rec.Body.String())
}
