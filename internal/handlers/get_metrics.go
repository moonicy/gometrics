package handlers

import (
	"fmt"
	"github.com/moonicy/gometrics/internal/utils"
	"net/http"
	"strings"
)

func (u *MetricsHandler) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	gotCounter, gotGauge := u.mem.GetMetrics()
	builder := strings.Builder{}
	for k, v := range gotCounter {
		builder.WriteString(fmt.Sprintf("%s: %d\n", k, v))
	}
	for k, v := range gotGauge {
		builder.WriteString(fmt.Sprintf("%s: %s\n", k, utils.FloatToString(v)))
	}
	_, err := res.Write([]byte(builder.String()))
	if err != nil {
		http.Error(res, "Internal Error", http.StatusInternalServerError)
	}
}
