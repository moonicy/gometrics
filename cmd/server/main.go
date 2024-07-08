package main

import (
	"github.com/moonicy/gometrics/internal/handlers"
	"github.com/moonicy/gometrics/internal/logger"
	"net/http"
)

func main() {
	addr := parseFlagRunAddr()
	sugar := logger.NewLogger()
	route := handlers.NewRoute(sugar)

	sugar.Infow(
		"Starting server",
		"addr", addr,
	)

	err := http.ListenAndServe(addr, route)
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
