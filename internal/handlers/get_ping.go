package handlers

import (
	"net/http"
)

// GetPing обрабатывает HTTP-запрос для проверки доступности сервера.
// Он выполняет операцию Ping через mh.pinger и возвращает соответствующий статус.
// В случае ошибки возвращает HTTP 500 Internal Server Error.
func (mh *MetricsHandler) GetPing(res http.ResponseWriter, _ *http.Request) {
	err := mh.pinger.Ping()
	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
	}
}
