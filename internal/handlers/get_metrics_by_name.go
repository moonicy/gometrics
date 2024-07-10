package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
	"github.com/moonicy/gometrics/pkg/floattostr"
	"net/http"
)

func (mh *MetricsHandler) GetMetricsByName(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, metrics.MName)
	tp := chi.URLParam(req, metrics.MType)

	if name == "" {
		http.Error(res, "Not found", http.StatusNotFound)
	}

	switch tp {
	case metrics.Gauge:
		value, err := mh.mem.GetGauge(req.Context(), name)
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
		value, err := mh.mem.GetCounter(req.Context(), name)
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
