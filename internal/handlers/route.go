package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moonicy/gometrics/pkg/middlewares"
	"go.uber.org/zap"
)

func NewRoute(mh *MetricsHandler, log zap.SugaredLogger) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Use(middlewares.GzipMiddleware)
		r.Use(middlewares.WithLogging(log))
		r.Get("/", mh.GetMetrics)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", mh.PostJSONMetricsByName)
			r.Get("/{type}/{name}", mh.GetMetricsByName)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", mh.UpdateJSONMetric)
			r.Post("/{type}/{name}/{value}", mh.UpdateMetric)
		})
		r.Route("/updates", func(r chi.Router) {
			r.Post("/", mh.UpdatesJSONMetrics)
		})
		r.Get("/ping", mh.GetPing)
	})

	return router
}
