package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// AgentConfig хранит информацию о конфгурации агента.
type AgentConfig struct {
	// Host - адрес эндпоинта HTTP-сервера.
	Host string `json:"address"`
	// ReportInterval - частота отправки метрик на сервер.
	ReportInterval time.Duration `json:"report_interval"`
	// PollInterval - частота опроса метрик из пакета runtime.
	PollInterval time.Duration `json:"poll_interval"`
	// HashKey - ключ для хеша.
	HashKey string
	// RateLimit - количество одновременно исходящих запросов на сервер.
	RateLimit int
	// CryptoKey - путь до файла с публичным ключом.
	CryptoKey string `json:"crypto_key"`
	// Config - путь до файла конфигурации.
	Config string
}

// NewAgentConfig создаёт и возвращает новый экземпляр AgentConfig, инициализированный с помощью флагов.
func NewAgentConfig() AgentConfig {
	ac := AgentConfig{}
	ac.parseFlag()
	return ac
}

func (ac *AgentConfig) parseFlag() {
	var scFlags AgentConfig
	var err error

	flag.StringVar(&scFlags.Host, "a", DefaultHost, "address and port to run server")
	flag.DurationVar(&scFlags.ReportInterval, "r", DefaultReportInterval*time.Second, "report interval")
	flag.DurationVar(&scFlags.PollInterval, "p", DefaultPollInterval*time.Second, "poll interval")
	flag.StringVar(&scFlags.HashKey, "k", DefaultHashKey, "hash key")
	flag.IntVar(&scFlags.RateLimit, "l", DefaultRateLimit, "rate limit")
	flag.StringVar(&scFlags.CryptoKey, "crypto-key", DefaultCryptoKeyAgent, "crypto key")
	flag.StringVar(&scFlags.Config, "c", "", "file config")
	flag.StringVar(&ac.Config, "config", "", "file config")
	flag.Parse()

	if scFlags.Config != "" {
		ac.Config = scFlags.Config
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		ac.Config = envConfig
	}
	if ac.Config != "" {
		file, errl := os.ReadFile(ac.Config)
		if errl != nil {
			log.Fatal(errl)
		}
		errl = json.Unmarshal(file, &ac)
		if errl != nil {
			log.Fatal(errl)
		}
	}

	if scFlags.Host != "" {
		ac.Host = scFlags.Host
	}
	if scFlags.ReportInterval > 0 {
		ac.ReportInterval = scFlags.ReportInterval
	}
	if scFlags.PollInterval > 0 {
		ac.PollInterval = scFlags.PollInterval
	}
	if scFlags.HashKey != "" {
		ac.HashKey = scFlags.HashKey
	}
	if scFlags.RateLimit > 0 {
		ac.RateLimit = scFlags.RateLimit
	}
	if scFlags.CryptoKey != "" {
		ac.CryptoKey = scFlags.CryptoKey
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		ac.Host = envRunAddr
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		ac.HashKey = envHashKey
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		str := strings.Trim(envReportInterval, "\"")
		if i, errStr := strconv.Atoi(str); errStr == nil {
			ac.ReportInterval = time.Duration(i)
		}

		if dur, errPD := time.ParseDuration(str); errPD == nil {
			ac.ReportInterval = dur
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		str := strings.Trim(envPollInterval, "\"")
		if i, errStr := strconv.Atoi(str); errStr == nil {
			ac.PollInterval = time.Duration(i)
		}

		if dur, errPD := time.ParseDuration(str); errPD == nil {
			ac.PollInterval = dur
		}
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		ac.RateLimit, err = strconv.Atoi(envRateLimit)
		if err != nil {
			log.Fatal("Invalid RATE_LIMIT")
		}
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		ac.CryptoKey = envCryptoKey
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		ac.Config = envConfig
	}
}
