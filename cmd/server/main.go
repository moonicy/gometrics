package main

import (
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

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var mem = NewMemStorage()

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func postUpdate(res http.ResponseWriter, req *http.Request) {
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
		mem.gauge[name] = valFloat
	case counter:
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			http.Error(res, "Value is not a valid int64", http.StatusBadRequest)
			return
		}
		mem.counter[name] += valInt
	default:
		http.Error(res, "Bad request", http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}/", postUpdate)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
