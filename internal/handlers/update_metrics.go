package handlers

import (
	"fmt"
	"github.com/moonicy/gometrics/internal/storage"
	"net/http"
	"strconv"
)

const (
	gauge   = "gauge"
	counter = "counter"
	mName   = "name"
	mValue  = "value"
	mType   = "type"
)

type UpdateMetrics struct {
	mem storage.MemoryStorage
}

func NewUpdateMetrics(mem storage.MemoryStorage) *UpdateMetrics {
	return &UpdateMetrics{mem}
}

func (u *UpdateMetrics) UpdateMetrics(res http.ResponseWriter, req *http.Request) {
	name := req.PathValue(mName)
	val := req.PathValue(mValue)
	tp := req.PathValue(mType)
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
