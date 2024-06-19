package main

import (
	"github.com/moonicy/gometrics/internal/handlers"
	"net/http"
)

func main() {
	route := handlers.Route()

	err := http.ListenAndServe(`:8080`, route)
	if err != nil {
		panic(err)
	}
}
