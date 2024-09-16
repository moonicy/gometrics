package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/middlewares"
)

func NewRoute(mh *MetricsHandler, log zap.SugaredLogger, cfg config.ServerConfig) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Use(middlewares.GzipMiddleware)
		r.Use(middlewares.WithLogging(log))
		r.Use(middlewares.SignCheckMiddleware(cfg.HashKey))
		r.Get("/", mh.GetMetrics)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", mh.GetMetricValueByNameJSON)
			r.Get("/{type}/{name}", mh.GetMetricValueByName)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", mh.PostMetricUpdateJSON)
			r.Post("/{type}/{name}/{value}", mh.PostMetricUpdate)
		})
		r.Route("/updates", func(r chi.Router) {
			r.Post("/", mh.PostMetricsUpdatesJSON)
		})
		r.Get("/ping", mh.GetPing)
	})

	return router
}
