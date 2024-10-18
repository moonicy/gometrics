package config

import (
	"os"
	"testing"
	"time"
)

func TestNewAgentConfig_Defaults(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	resetFlags()

	os.Args = []string{"cmd"}

	ac := NewAgentConfig()

	if ac.Host != DefaultHost {
		t.Errorf("Expected Host to be '%s', got '%s'", DefaultHost, ac.Host)
	}
	if ac.ReportInterval != DefaultReportInterval*time.Second {
		t.Errorf("Expected ReportInterval to be %v, got %v", DefaultReportInterval*time.Second, ac.ReportInterval)
	}
	if ac.PollInterval != DefaultPollInterval*time.Second {
		t.Errorf("Expected PollInterval to be %v, got %v", DefaultPollInterval*time.Second, ac.PollInterval)
	}
	if ac.HashKey != DefaultHashKey {
		t.Errorf("Expected HashKey to be '%s', got '%s'", DefaultHashKey, ac.HashKey)
	}
	if ac.RateLimit != DefaultRateLimit {
		t.Errorf("Expected RateLimit to be %d, got %d", DefaultRateLimit, ac.RateLimit)
	}
	if ac.CryptoKey != DefaultCryptoKeyAgent {
		t.Errorf("Expected CryptoKey to be '%s', got '%s'", DefaultCryptoKeyAgent, ac.CryptoKey)
	}
	if ac.Config != "" {
		t.Errorf("Expected Config to be empty, got '%s'", ac.Config)
	}
}

func TestNewAgentConfig_Flags(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	resetFlags()

	os.Args = []string{
		"cmd",
		"-a", "localhost:8081",
		"-r", "15s",
		"-p", "10s",
		"-k", "secretkey",
		"-l", "5",
		"-crypto-key", "/path/to/crypto.key",
	}

	ac := NewAgentConfig()

	if ac.Host != "localhost:8081" {
		t.Errorf("Expected Host to be 'localhost:8081', got '%s'", ac.Host)
	}
	if ac.ReportInterval != 15*time.Second {
		t.Errorf("Expected ReportInterval to be 15s, got %v", ac.ReportInterval)
	}
	if ac.PollInterval != 10*time.Second {
		t.Errorf("Expected PollInterval to be 10s, got %v", ac.PollInterval)
	}
	if ac.HashKey != "secretkey" {
		t.Errorf("Expected HashKey to be 'secretkey', got '%s'", ac.HashKey)
	}
	if ac.RateLimit != 5 {
		t.Errorf("Expected RateLimit to be 5, got %d", ac.RateLimit)
	}
	if ac.CryptoKey != "/path/to/crypto.key" {
		t.Errorf("Expected CryptoKey to be '/path/to/crypto.key', got '%s'", ac.CryptoKey)
	}
}

func TestNewAgentConfig_EnvVars(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	resetFlags()

	err := os.Setenv("ADDRESS", "localhost:8082")
	if err != nil {
		t.Errorf("Failed to set environment variable ADDRESS: %v", err)
	}
	err = os.Setenv("REPORT_INTERVAL", "20s")
	if err != nil {
		t.Errorf("Failed to set environment variable REPORT_INTERVAL: %v", err)
	}
	err = os.Setenv("POLL_INTERVAL", "15s")
	if err != nil {
		t.Errorf("Failed to set environment variable POLL_INTERVAL: %v", err)
	}
	err = os.Setenv("KEY", "envsecretkey")
	if err != nil {
		t.Errorf("Failed to set environment variable KEY: %v", err)
	}
	err = os.Setenv("RATE_LIMIT", "10")
	if err != nil {
		t.Errorf("Failed to set environment variable RATE_LIMIT: %v", err)
	}
	err = os.Setenv("CRYPTO_KEY", "/env/path/to/crypto.key")
	if err != nil {
		t.Errorf("Failed to set environment variable CRYPTO_KEY: %v", err)
	}
	err = os.Setenv("CONFIG", "")
	if err != nil {
		t.Errorf("Failed to set environment variable CONFIG: %v", err)
	}

	ac := NewAgentConfig()

	if ac.Host != "localhost:8082" {
		t.Errorf("Expected Host to be 'localhost:8082', got '%s'", ac.Host)
	}
	if ac.ReportInterval != 20*time.Second {
		t.Errorf("Expected ReportInterval to be 20s, got %v", ac.ReportInterval)
	}
	if ac.PollInterval != 15*time.Second {
		t.Errorf("Expected PollInterval to be 15s, got %v", ac.PollInterval)
	}
	if ac.HashKey != "envsecretkey" {
		t.Errorf("Expected HashKey to be 'envsecretkey', got '%s'", ac.HashKey)
	}
	if ac.RateLimit != 10 {
		t.Errorf("Expected RateLimit to be 10, got %d", ac.RateLimit)
	}
	if ac.CryptoKey != "/env/path/to/crypto.key" {
		t.Errorf("Expected CryptoKey to be '/env/path/to/crypto.key', got '%s'", ac.CryptoKey)
	}
	if ac.Config != "" {
		t.Errorf("Expected Config to be '', got '%s'", ac.Config)
	}
}
