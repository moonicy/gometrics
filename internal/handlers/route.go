package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
	"github.com/moonicy/gometrics/internal/middlewares"
	"github.com/moonicy/gometrics/internal/storage"
	"go.uber.org/zap"
)

func NewRoute(ctx context.Context, log zap.SugaredLogger, cfg config.ServerConfig) *chi.Mux {
	cm := file.NewConsumer(cfg.FileStoragePath)
	pr := file.NewProducer(cfg.FileStoragePath)
	st := storage.NewFileStorage(ctx, cfg, cm, pr)
	metricsHandler := NewMetricsHandler(st)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Use(middlewares.GzipMiddleware)
		r.Use(middlewares.WithLogging(log))
		r.Get("/", metricsHandler.GetMetrics)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", metricsHandler.PostJSONMetricsByName)
			r.Get("/{type}/{name}", metricsHandler.GetMetricsByName)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", metricsHandler.UpdateJSONMetrics)
			r.Post("/{type}/{name}/{value}", metricsHandler.UpdateMetrics)
		})
	})

	return router
}
