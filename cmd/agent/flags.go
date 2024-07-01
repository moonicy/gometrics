package main

import (
	"flag"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/http"
	"os"
)

//var flagRunAddr string
//var flagReportInterval string
//var flagPollInterval string

func parseFlags() config.AgentConfig {
	cfg := config.AgentConfig{}
	flag.StringVar(&cfg.Host, "a", http.DefaultHost, "address and port to run server")
	flag.StringVar(&cfg.ReportInterval, "r", "10", "report interval")
	flag.StringVar(&cfg.PollInterval, "p", "2", "poll interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Host = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		cfg.ReportInterval = envReportInterval
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		cfg.PollInterval = envPollInterval
	}
	return cfg
}
