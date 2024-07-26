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
	HashKey        string
	RateLimit      int
}

func NewAgentConfig() AgentConfig {
	ac := AgentConfig{}
	ac.parseFlag()
	return ac
}

func (ac *AgentConfig) parseFlag() {
	var err error
	flag.StringVar(&ac.Host, "a", DefaultHost, "address and port to run server")
	flag.IntVar(&ac.ReportInterval, "r", 20, "report interval")
	flag.IntVar(&ac.PollInterval, "p", 2, "poll interval")
	flag.StringVar(&ac.HashKey, "k", "", "hash key")
	flag.IntVar(&ac.RateLimit, "l", 0, "rate limit")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		ac.Host = envRunAddr
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		ac.HashKey = envHashKey
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
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		ac.RateLimit, err = strconv.Atoi(envRateLimit)
		if err != nil {
			log.Fatal("Invalid RATE_LIMIT")
		}
	}
}
