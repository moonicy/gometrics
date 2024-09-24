package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
	"github.com/moonicy/gometrics/pkg/floattostr"
)

// GetMetricValueByName обрабатывает HTTP-запрос для получения значения метрики по её имени и типу.
// Он извлекает параметры из URL и возвращает значение метрики клиенту.
// В случае ошибки возвращает соответствующий HTTP-статус и сообщение об ошибке.
func (mh *MetricsHandler) GetMetricValueByName(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, metrics.MName)
	tp := chi.URLParam(req, metrics.MType)
	if mh.logger != nil {
		mh.logger.Infoln(name, tp)
	}

	if name == "" {
		http.Error(res, "Not found", http.StatusNotFound)
	}

	switch tp {
	case metrics.Gauge:
		value, err := mh.storage.GetGauge(req.Context(), name)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(res, "Not found", http.StatusNotFound)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = res.Write([]byte(floattostr.FloatToString(value)))
		if err != nil {
			http.Error(res, "Internal Error", http.StatusInternalServerError)
		}
	case metrics.Counter:
		value, err := mh.storage.GetCounter(req.Context(), name)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(res, "Not found", http.StatusNotFound)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = res.Write([]byte(fmt.Sprintf("%d", value)))
		if err != nil {
			http.Error(res, "Internal Error", http.StatusInternalServerError)
		}
	default:
		http.Error(res, "Bad request", http.StatusBadRequest)
	}
}
