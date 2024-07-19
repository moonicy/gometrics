package handlers

import (
	"net/http"
)

func (mh *MetricsHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	err := mh.pinger.Ping()
	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
	}
}
