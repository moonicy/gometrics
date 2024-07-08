package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/metrics"
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
		value, ok := mh.mem.GetGauge(name)
		if !ok {
			http.Error(res, "Not found", http.StatusNotFound)
			return
		}
		_, err := res.Write([]byte(floattostr.FloatToString(value)))
		if err != nil {
			http.Error(res, "Internal Error", http.StatusInternalServerError)
		}
	case metrics.Counter:
		value, ok := mh.mem.GetCounter(name)
		if !ok {
			http.Error(res, "Not found", http.StatusNotFound)
			return
		}
		_, err := res.Write([]byte(fmt.Sprintf("%d", value)))
		if err != nil {
			http.Error(res, "Internal Error", http.StatusInternalServerError)
		}
	default:
		http.Error(res, "Bad request", http.StatusBadRequest)
	}
}
