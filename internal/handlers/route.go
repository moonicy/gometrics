package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/middlewares"
	"github.com/moonicy/gometrics/internal/storage"
	"go.uber.org/zap"
)

func NewRoute(log zap.SugaredLogger) *chi.Mux {
	mem := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(mem)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/", middlewares.WithLogging(log, metricsHandler.GetMetrics))
		r.Get("/value/{type}/{name}", middlewares.WithLogging(log, metricsHandler.GetMetricsByName))
		r.Post("/update/{type}/{name}/{value}", middlewares.WithLogging(log, metricsHandler.UpdateMetrics))
	})

	return router
}
