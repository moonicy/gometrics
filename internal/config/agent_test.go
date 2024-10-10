package config

import (
	"os"
	"testing"
)

func TestNewAgentConfig_Flags(t *testing.T) {
	resetFlags()

	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:9090",
		"-r", "120",
		"-p", "20",
		"-k", "testhashkey",
		"-l", "10",
	}

	ac := NewAgentConfig()

	if ac.Host != "127.0.0.1:9090" {
		t.Errorf("Expected Host to be %s, got %s", "127.0.0.1:9090", ac.Host)
	}
	if ac.ReportInterval != 120 {
		t.Errorf("Expected ReportInterval to be %d, got %d", 120, ac.ReportInterval)
	}
	if ac.PollInterval != 20 {
		t.Errorf("Expected PollInterval to be %d, got %d", 20, ac.PollInterval)
	}
	if ac.HashKey != "testhashkey" {
		t.Errorf("Expected HashKey to be %s, got %s", "testhashkey", ac.HashKey)
	}
	if ac.RateLimit != 10 {
		t.Errorf("Expected RateLimit to be %d, got %d", 10, ac.RateLimit)
	}
}

func TestNewAgentConfig_EnvVars(t *testing.T) {
	resetFlags()

	os.Setenv("ADDRESS", "192.168.1.1:8081")
	os.Setenv("KEY", "envhashkey")
	os.Setenv("REPORT_INTERVAL", "180")
	os.Setenv("POLL_INTERVAL", "30")
	os.Setenv("RATE_LIMIT", "15")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("KEY")
		os.Unsetenv("REPORT_INTERVAL")
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("RATE_LIMIT")
	}()

	ac := NewAgentConfig()

	if ac.Host != "192.168.1.1:8081" {
		t.Errorf("Expected Host to be %s, got %s", "192.168.1.1:8081", ac.Host)
	}
	if ac.ReportInterval != 180 {
		t.Errorf("Expected ReportInterval to be %d, got %d", 180, ac.ReportInterval)
	}
	if ac.PollInterval != 30 {
		t.Errorf("Expected PollInterval to be %d, got %d", 30, ac.PollInterval)
	}
	if ac.HashKey != "envhashkey" {
		t.Errorf("Expected HashKey to be %s, got %s", "envhashkey", ac.HashKey)
	}
	if ac.RateLimit != 15 {
		t.Errorf("Expected RateLimit to be %d, got %d", 15, ac.RateLimit)
	}
}
