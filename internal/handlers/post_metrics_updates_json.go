package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/moonicy/gometrics/internal/metrics"
)

func (mh *MetricsHandler) PostMetricsUpdatesJSON(res http.ResponseWriter, req *http.Request) {
	var mt []metrics.Metric

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if mh.logger != nil {
		mh.logger.Infoln(string(body))
	}
	if err = json.Unmarshal(body, &mt); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if len(mt) == 0 {
		http.Error(res, "no metrics found", http.StatusBadRequest)
	}
	for _, m := range mt {
		if err = m.Validate(); err != nil {
			if errors.Is(err, metrics.ErrNotFound) {
				http.Error(res, err.Error(), http.StatusNotFound)
			}
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
	}
	mtGauge := make(map[string]float64)
	mtCounter := make(map[string]int64)
	for _, m := range mt {
		if m.MType == metrics.Gauge {
			mtGauge[m.ID] = *m.Value
		} else {
			mtCounter[m.ID] += *m.Delta
		}
	}
	err = mh.storage.SetMetrics(req.Context(), mtCounter, mtGauge)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
