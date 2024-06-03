package main

import (
	"github.com/moonicy/gometrics/internal/handlers"
	"net/http"
)

func main() {
	mux := handlers.Route()

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
