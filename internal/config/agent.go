package config

import (
	"flag"
	"log"
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
	var err error
	flag.StringVar(&ac.Host, "a", DefaultHost, "address and port to run server")
	flag.IntVar(&ac.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&ac.PollInterval, "p", 2, "poll interval")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		ac.Host = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		ac.ReportInterval, err = strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatal("Invalid REPORT_INTERVAL")
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		ac.PollInterval, err = strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatal("Invalid POLL_INTERVAL")
		}
	}
}