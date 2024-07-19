package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func (mh *MetricsHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	err := mh.db.Ping()
	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
	}
}
