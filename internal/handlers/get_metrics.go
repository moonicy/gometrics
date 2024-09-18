package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/moonicy/gometrics/pkg/floattostr"
)

// GetMetrics обрабатывает HTTP-запрос для получения значения всех метрик.
// В случае ошибки возвращает соответствующий HTTP-статус и сообщение об ошибке.
func (mh *MetricsHandler) GetMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	gotCounter, gotGauge, err := mh.storage.GetMetrics(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	builder := strings.Builder{}
	for k, v := range gotCounter {
		builder.WriteString(fmt.Sprintf("%s: %d\n", k, v))
	}
	for k, v := range gotGauge {
		builder.WriteString(fmt.Sprintf("%s: %s\n", k, floattostr.FloatToString(v)))
	}
	_, err = res.Write([]byte(builder.String()))
	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
	}
}
