package handlers

import (
	"encoding/json"
	"errors"
	"github.com/moonicy/gometrics/internal/metrics"
	"io"
	"log"
	"net/http"
)

func (u *MetricsHandler) PostJSONMetricsByName(res http.ResponseWriter, req *http.Request) {
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
		value, ok := u.mem.GetGauge(mt.ID)
		if !ok {
			http.Error(res, "Not found", http.StatusNotFound)
			return
		}
		resBody := metrics.Metrics{MetricName: metrics.MetricName{ID: mt.ID, MType: mt.MType}, Value: &value}
		out, err := json.Marshal(resBody)
		if err != nil {
			log.Fatal(err)
		}
		res.Write(out)
	case metrics.Counter:
		delta, ok := u.mem.GetCounter(mt.ID)
		if !ok {
			http.Error(res, "Not found", http.StatusNotFound)
			return
		}
		resBody := metrics.Metrics{MetricName: metrics.MetricName{ID: mt.ID, MType: mt.MType}, Delta: &delta}
		out, err := json.Marshal(resBody)
		if err != nil {
			log.Fatal(err)
		}
		res.Write(out)
	}
}
