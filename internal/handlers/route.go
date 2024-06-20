package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/storage"
)

func RouteNew() *chi.Mux {
	mem := storage.NewMemStorage()
	metricsHandler := NewMetricsHandler(mem)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/", metricsHandler.GetMetrics)
		r.Get("/value/{type}/{name}", metricsHandler.GetMetricsByName)
		r.Post("/update/{type}/{name}/{value}", metricsHandler.UpdateMetrics)
	})

	return router
}
