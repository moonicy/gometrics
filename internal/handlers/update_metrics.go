package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (u *MetricsHandler) UpdateMetrics(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, mName)
	val := chi.URLParam(req, mValue)
	tp := chi.URLParam(req, mType)

	if name == "" {
		http.Error(res, "Not found", http.StatusNotFound)
	}

	switch tp {
	case gauge:
		valFloat, err := strconv.ParseFloat(val, 64)
		if err != nil {
			http.Error(res, "Value is not a valid float64", http.StatusBadRequest)
			return
		}
		u.mem.SetGauge(name, valFloat)
	case counter:
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(res, "Value is not a valid int64", http.StatusBadRequest)
			return
		}
		u.mem.AddCounter(name, valInt)
	default:
		http.Error(res, "Bad request", http.StatusBadRequest)
	}
	fmt.Printf("%s\t%s\t%s\n", name, val, tp)
}
