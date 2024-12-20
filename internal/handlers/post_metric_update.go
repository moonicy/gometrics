package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/moonicy/gometrics/internal/metrics"
)

// PostMetricUpdate обрабатывает HTTP-запрос для обновления значения метрики.
// Получает имя метрики, тип и значение из параметров URL и обновляет хранилище метрик.
// В случае ошибки возвращает соответствующий HTTP-статус и сообщение об ошибке.
func (mh *MetricsHandler) PostMetricUpdate(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, metrics.MName)
	val := chi.URLParam(req, metrics.MValue)
	tp := chi.URLParam(req, metrics.MType)

	if name == "" {
		http.Error(res, "Not found", http.StatusNotFound)
	}

	switch tp {
	case metrics.Gauge:
		valFloat, err := strconv.ParseFloat(val, 64)
		if err != nil {
			http.Error(res, "Value is not a valid float64", http.StatusBadRequest)
			return
		}
		err = mh.storage.SetGauge(req.Context(), name, valFloat)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	case metrics.Counter:
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(res, "Value is not a valid int64", http.StatusBadRequest)
			return
		}
		err = mh.storage.AddCounter(req.Context(), name, valInt)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(res, "Bad request", http.StatusBadRequest)
	}
	fmt.Printf("%s\t%s\t%s\n", name, val, tp)
}
