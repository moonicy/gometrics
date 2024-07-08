package main

import (
	"flag"
	"github.com/moonicy/gometrics/internal/http"
	"os"
)

var flagRunAddr string
var flagReportInterval string
var flagPollInterval string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", http.DefaultHost, "address and port to run server")
	flag.StringVar(&flagReportInterval, "r", "10", "report interval")
	flag.StringVar(&flagPollInterval, "p", "2", "poll interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		flagReportInterval = envReportInterval
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		flagPollInterval = envPollInterval
	}
}
