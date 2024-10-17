package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseSize += size
	return size, err
}

// WithLogging возвращает middleware, который логирует информацию о каждом HTTP-запросе и ответе.
// Он записывает URI запроса, метод, длительность обработки, статусный код и размер ответа.
func WithLogging(sugar *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{ResponseWriter: res, statusCode: http.StatusOK}

			handler.ServeHTTP(lrw, req)

			duration := time.Since(start)

			sugar.Infoln(
				"uri", req.RequestURI,
				"method", req.Method,
				"duration", duration,
			)

			sugar.Infoln(
				"status", lrw.statusCode,
				"responseSize", lrw.responseSize,
			)
		})
	}

}
