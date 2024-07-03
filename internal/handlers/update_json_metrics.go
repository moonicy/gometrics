package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moonicy/gometrics/internal/metrics"
	"io"
	"log"
	"net/http"
)

func (mh *MetricsHandler) UpdateJSONMetrics(res http.ResponseWriter, req *http.Request) {
	var mt metrics.Metric

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if err = json.Unmarshal(body, &mt); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	if err = mt.Validate(); err != nil {
		if errors.Is(err, metrics.ErrNotFound) {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	var value *float64
	var delta *int64
	switch mt.MType {
	case metrics.Gauge:
		mh.mem.SetGauge(mt.ID, *mt.Value)
		gv, ok := mh.mem.GetGauge(mt.ID)
		if ok {
			value = &gv
		}
		fmt.Printf("%s\t%s\t%f\n", mt.ID, mt.MType, *mt.Value)
	case metrics.Counter:
		mh.mem.AddCounter(mt.ID, *mt.Delta)
		cv, ok := mh.mem.GetCounter(mt.ID)
		if ok {
			delta = &cv
		}
		fmt.Printf("%s\t%s\t%d\n", mt.ID, mt.MType, *mt.Delta)
	}

	resBody := metrics.Metric{MetricName: metrics.MetricName{ID: mt.ID, MType: mt.MType}, Value: value, Delta: delta}
	out, err := json.Marshal(resBody)
	if err != nil {
		log.Fatal(err)
	}
	res.Write(out)
}
