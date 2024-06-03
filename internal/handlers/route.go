package handlers

import (
	"github.com/moonicy/gometrics/internal/storage"
	"net/http"
)

func Route() *http.ServeMux {
	mem := storage.NewMemStorage()
	updateMetrics := NewUpdateMetrics(mem)
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}/", updateMetrics.UpdateMetrics)
	return mux
}
