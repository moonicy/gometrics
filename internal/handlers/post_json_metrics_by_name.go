package handlers

import (
	"encoding/json"
	"errors"
	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
	"io"
	"log"
	"net/http"
)

func (mh *MetricsHandler) PostJSONMetricsByName(res http.ResponseWriter, req *http.Request) {
	var mt metrics.MetricName

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &mt); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = mt.Validate(); err != nil {
		if errors.Is(err, metrics.ErrNotFound) {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	switch mt.MType {
	case metrics.Gauge:
		value, err := mh.mem.GetGauge(req.Context(), mt.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(res, "Not found", http.StatusNotFound)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		resBody := metrics.Metric{MetricName: metrics.MetricName{ID: mt.ID, MType: mt.MType}, Value: &value}
		out, err := json.Marshal(resBody)
		if err != nil {
			log.Fatal(err)
		}
		res.Write(out)
	case metrics.Counter:
		delta, err := mh.mem.GetCounter(req.Context(), mt.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(res, "Not found", http.StatusNotFound)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		resBody := metrics.Metric{MetricName: metrics.MetricName{ID: mt.ID, MType: mt.MType}, Delta: &delta}
		out, err := json.Marshal(resBody)
		if err != nil {
			log.Fatal(err)
		}
		res.Write(out)
	}
}
