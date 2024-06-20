package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func WithLogging(sugar zap.SugaredLogger, h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)

	}
	return logFn
}
