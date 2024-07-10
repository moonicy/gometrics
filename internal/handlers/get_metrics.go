package handlers

import (
	"fmt"
	"github.com/moonicy/gometrics/pkg/floattostr"
	"net/http"
	"strings"
)

func (mh *MetricsHandler) GetMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	gotCounter, gotGauge, err := mh.mem.GetMetrics(req.Context())
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
