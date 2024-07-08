package config

import (
	"flag"
	"os"
	"strconv"
)

type AgentConfig struct {
	Host           string
	ReportInterval int
	PollInterval   int
}

func NewAgentConfig() AgentConfig {
	ac := AgentConfig{}
	ac.parseFlag()
	return ac
}

func (ac *AgentConfig) parseFlag() {
	flag.StringVar(&ac.Host, "a", DefaultHost, "address and port to run server")
	flag.IntVar(&ac.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&ac.PollInterval, "p", 2, "poll interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		ac.Host = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		ac.ReportInterval, _ = strconv.Atoi(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		ac.PollInterval, _ = strconv.Atoi(envPollInterval)
	}
}
