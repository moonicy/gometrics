package config

import (
	"flag"
	"os"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestNewServerConfig_WithEnvVars(t *testing.T) {
	resetFlags()
	os.Clearenv()

	os.Setenv("ADDRESS", "192.168.1.1:8081")
	os.Setenv("STORE_INTERVAL", "900")
	os.Setenv("FILE_STORAGE_PATH", "/var/data")
	os.Setenv("DATABASE_DSN", "user:pass@/prod_db")
	os.Setenv("KEY", "envsecret")
	os.Setenv("RESTORE", "false")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("KEY")
		os.Unsetenv("RESTORE")
	}()

	sc := NewServerConfig()

	if sc.Host != "192.168.1.1:8081" {
		t.Errorf("Expected Host to be 192.168.1.1:8081, got %s", sc.Host)
	}
	if sc.StoreInternal != 900 {
		t.Errorf("Expected StoreInternal to be 900, got %d", sc.StoreInternal)
	}
	if sc.FileStoragePath != "/var/data" {
		t.Errorf("Expected FileStoragePath to be /var/data, got %s", sc.FileStoragePath)
	}
	if sc.Restore != false {
		t.Errorf("Expected Restore to be false, got %v", sc.Restore)
	}
	if sc.DatabaseDsn != "user:pass@/prod_db" {
		t.Errorf("Expected DatabaseDsn to be user:pass@/prod_db, got %s", sc.DatabaseDsn)
	}
	if sc.HashKey != "envsecret" {
		t.Errorf("Expected HashKey to be envsecret, got %s", sc.HashKey)
	}
}

func TestNewServerConfig_FlagsAndEnvVars(t *testing.T) {
	resetFlags()
	os.Clearenv()

	os.Setenv("ADDRESS", "192.168.1.1:8081")
	os.Setenv("STORE_INTERVAL", "900")

	os.Args = []string{
		"cmd",
		"-a", "127.0.0.1:9090",
		"-i", "600",
	}

	sc := NewServerConfig()

	if sc.Host != "192.168.1.1:8081" {
		t.Errorf("Expected Host to be 192.168.1.1:8081, got %s", sc.Host)
	}
	if sc.StoreInternal != 900 {
		t.Errorf("Expected StoreInternal to be 900, got %d", sc.StoreInternal)
	}
	if sc.FileStoragePath != "" {
		t.Errorf("Expected FileStoragePath to be empty, got %s", sc.FileStoragePath)
	}
	if sc.Restore != true {
		t.Errorf("Expected Restore to be true, got %v", sc.Restore)
	}
	if sc.DatabaseDsn != "" {
		t.Errorf("Expected DatabaseDsn to be empty, got %s", sc.DatabaseDsn)
	}
	if sc.HashKey != "" {
		t.Errorf("Expected HashKey to be empty, got %s", sc.HashKey)
	}
}

func TestNewServerConfig_InvalidEnvVars(t *testing.T) {
	resetFlags()
	os.Clearenv()

	os.Setenv("STORE_INTERVAL", "invalid")
	os.Setenv("RESTORE", "invalid")
	defer func() {
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("RESTORE")
	}()

	sc := NewServerConfig()

	if sc.StoreInternal != 0 {
		t.Errorf("Expected StoreInternal to be 0 due to invalid env var, got %d", sc.StoreInternal)
	}
	if sc.Restore != true {
		t.Errorf("Expected Restore to remain true due to invalid env var, got %v", sc.Restore)
	}
}
