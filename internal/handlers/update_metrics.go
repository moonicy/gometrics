package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moonicy/gometrics/internal/metrics"
	"io"
	"net/http"
)

func (u *MetricsHandler) UpdateMetrics(res http.ResponseWriter, req *http.Request) {

	var mt metrics.Metrics

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
		u.mem.SetGauge(mt.ID, *mt.Value)
		fmt.Printf("%s\t%s\t%f\n", mt.ID, mt.MType, *mt.Value)
	case metrics.Counter:
		u.mem.AddCounter(mt.ID, *mt.Delta)
		fmt.Printf("%s\t%s\t%d\n", mt.ID, mt.MType, *mt.Delta)
	}
}
