package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// AgentConfig хранит информацию о конфгурации агента.
type AgentConfig struct {
	// Host - адрес эндпоинта HTTP-сервера.
	Host string
	// ReportInterval - частота отправки метрик на сервер.
	ReportInterval int
	// PollInterval - частота опроса метрик из пакета runtime.
	PollInterval int
	// HashKey - ключ для хеша.
	HashKey string
	// RateLimit - количество одновременно исходящих запросов на сервер.
	RateLimit int
}

// NewAgentConfig создаёт и возвращает новый экземпляр AgentConfig, инициализированный с помощью флагов.
func NewAgentConfig() AgentConfig {
	ac := AgentConfig{}
	ac.parseFlag()
	return ac
}

func (ac *AgentConfig) parseFlag() {
	var err error
	flag.StringVar(&ac.Host, "a", DefaultHost, "address and port to run server")
	flag.IntVar(&ac.ReportInterval, "r", DefaultReportInterval, "report interval")
	flag.IntVar(&ac.PollInterval, "p", DefaultPollInterval, "poll interval")
	flag.StringVar(&ac.HashKey, "k", DefaultHashKey, "hash key")
	flag.IntVar(&ac.RateLimit, "l", DefaultRateLimit, "rate limit")
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
