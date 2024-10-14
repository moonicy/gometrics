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

	err := os.Setenv("ADDRESS", "192.168.1.1:8081")
	if err != nil {
		t.Errorf("Failed to set environment variable ADDRESS: %v", err)
	}
	err = os.Setenv("STORE_INTERVAL", "900")
	if err != nil {
		t.Errorf("Failed to set environment variable STORE_INTERVAL: %v", err)
	}
	err = os.Setenv("FILE_STORAGE_PATH", "/var/data")
	if err != nil {
		t.Errorf("Failed to set environment variable FILE_STORAGE_PATH: %v", err)
	}
	err = os.Setenv("DATABASE_DSN", "user:pass@/prod_db")
	if err != nil {
		t.Errorf("Failed to set environment variable DATABASE_DSN: %v", err)
	}
	err = os.Setenv("KEY", "envsecret")
	if err != nil {
		t.Errorf("Failed to set environment variable KEY: %v", err)
	}
	err = os.Setenv("RESTORE", "false")
	if err != nil {
		t.Errorf("Failed to set environment variable RESTORE: %v", err)
	}
	err = os.Setenv("CRYPTO_KEY", "str")
	if err != nil {
		t.Errorf("Failed to set environment variable CRYPTO_KEY: %v", err)
	}
	defer func() {
		err = os.Unsetenv("ADDRESS")
		if err != nil {
			t.Errorf("Failed to delete environment variable ADDRESS: %v", err)
		}
		err = os.Unsetenv("STORE_INTERVAL")
		if err != nil {
			t.Errorf("Failed to delete environment variable STORE_INTERVAL: %v", err)
		}
		err = os.Unsetenv("FILE_STORAGE_PATH")
		if err != nil {
			t.Errorf("Failed to delete environment variable FILE_STORAGE_PATH: %v", err)
		}
		err = os.Unsetenv("DATABASE_DSN")
		if err != nil {
			t.Errorf("Failed to delete environment variable DATABASE_DSN: %v", err)
		}
		err = os.Unsetenv("KEY")
		if err != nil {
			t.Errorf("Failed to delete environment variable KEY: %v", err)
		}
		err = os.Unsetenv("RESTORE")
		if err != nil {
			t.Errorf("Failed to delete environment variable RESTORE: %v", err)
		}
		err = os.Unsetenv("CRYPTO_KEY")
		if err != nil {
			t.Errorf("Failed to delete environment variable CRYPTO_KEY: %v", err)
		}
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
	if sc.CryptoKey != "str" {
		t.Errorf("Expected CryptoKey to be str, got %s", sc.CryptoKey)
	}
}

func TestNewServerConfig_FlagsAndEnvVars(t *testing.T) {
	resetFlags()
	os.Clearenv()

	err := os.Setenv("ADDRESS", "192.168.1.1:8081")
	if err != nil {
		t.Errorf("Failed to set environment variable ADDRESS: %v", err)
	}
	err = os.Setenv("STORE_INTERVAL", "900")
	if err != nil {
		t.Errorf("Failed to set environment variable STORE_INTERVAL: %v", err)
	}

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
	if sc.CryptoKey != "keys/private.pem" {
		t.Errorf("Expected CryptoKey to be keys/private.pem, got %s", sc.CryptoKey)
	}
}

func TestNewServerConfig_InvalidEnvVars(t *testing.T) {
	resetFlags()
	os.Clearenv()

	err := os.Setenv("STORE_INTERVAL", "invalid")
	if err != nil {
		t.Errorf("Failed to set environment variable STORE_INTERVAL: %v", err)
	}
	err = os.Setenv("RESTORE", "invalid")
	if err != nil {
		t.Errorf("Failed to set environment variable RESTORE: %v", err)
	}
	defer func() {
		err = os.Unsetenv("STORE_INTERVAL")
		if err != nil {
			t.Errorf("Failed to delete environment variable STORE_INTERVAL: %v", err)
		}
		err = os.Unsetenv("RESTORE")
		if err != nil {
			t.Errorf("Failed to delete environment variable RESTORE: %v", err)
		}
	}()

	sc := NewServerConfig()

	if sc.StoreInternal != 0 {
		t.Errorf("Expected StoreInternal to be 0 due to invalid env var, got %d", sc.StoreInternal)
	}
	if sc.Restore != true {
		t.Errorf("Expected Restore to remain true due to invalid env var, got %v", sc.Restore)
	}
}
