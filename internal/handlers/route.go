package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/middlewares"
)

type MetricsHandlers interface {
	GetMetrics(res http.ResponseWriter, req *http.Request)
	GetMetricValueByNameJSON(res http.ResponseWriter, req *http.Request)
	GetMetricValueByName(res http.ResponseWriter, req *http.Request)
	PostMetricUpdateJSON(res http.ResponseWriter, req *http.Request)
	PostMetricUpdate(res http.ResponseWriter, req *http.Request)
	PostMetricsUpdatesJSON(res http.ResponseWriter, req *http.Request)
	GetPing(res http.ResponseWriter, req *http.Request)
}

// NewRoute создаёт и настраивает новый маршрутизатор chi.Mux с необходимыми маршрутами и middleware.
// Он принимает MetricsHandler для обработки HTTP-запросов метрик, логгер и конфигурацию сервера.
func NewRoute(mh MetricsHandlers, log *zap.SugaredLogger, cfg config.ServerConfig) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Use(middlewares.CryptMiddleware("", cfg.CryptoKey))
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
			r.Use(middlewares.IPCheckMiddleware(cfg.TrustedSubnet))
			r.Post("/", mh.PostMetricsUpdatesJSON)
		})
		r.Get("/ping", mh.GetPing)
	})

	return router
}
