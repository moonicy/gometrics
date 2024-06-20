package main

import (
	"github.com/moonicy/gometrics/internal/handlers"
	"net/http"
)

func main() {
	parseFlags()

	route := handlers.RouteNew()

	err := http.ListenAndServe(flagRunAddr, route)
	if err != nil {
		panic(err)
	}
}
