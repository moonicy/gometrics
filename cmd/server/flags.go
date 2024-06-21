package main

import (
	"flag"
	"github.com/moonicy/gometrics/internal/http"
	"os"
)

func parseFlagRunAddr() string {
	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", http.DefaultHost, "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	return flagRunAddr
}
