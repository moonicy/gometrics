package config

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewServerConfig_Defaults(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd"}
	resetFlags()

	sc := NewServerConfig()

	if sc.Host != DefaultHost {
		t.Errorf("Expected Host to be '%s', got '%s'", DefaultHost, sc.Host)
	}
	if sc.StoreInterval != 300*time.Second {
		t.Errorf("Expected StoreInterval to be 300s, got %v", sc.StoreInterval)
	}
	if sc.FileStoragePath != "" {
		t.Errorf("Expected FileStoragePath to be empty, got '%s'", sc.FileStoragePath)
	}
	if sc.DatabaseDsn != "" {
		t.Errorf("Expected DatabaseDsn to be empty, got '%s'", sc.DatabaseDsn)
	}
	if sc.HashKey != "" {
		t.Errorf("Expected HashKey to be empty, got '%s'", sc.HashKey)
	}
	if sc.CryptoKey != DefaultCryptoKeyServer {
		t.Errorf("Expected CryptoKey to be '%s', got '%s'", DefaultCryptoKeyServer, sc.CryptoKey)
	}
	if sc.Config != "" {
		t.Errorf("Expected Config to be empty, got '%s'", sc.Config)
	}
}

func TestNewServerConfig_Flags(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	resetFlags()

	os.Args = []string{
		"cmd",
		"-a", "localhost:8081",
		"-i", "500s",
		"-f", "/tmp/storage.json",
		"-r", "false",
		"-d", "postgres://user:pass@localhost/db",
		"-k", "secret",
		"-crypto-key", "/path/to/key",
	}

	sc := NewServerConfig()

	if sc.Host != "localhost:8081" {
		t.Errorf("Expected Host to be 'localhost:8081', got '%s'", sc.Host)
	}
	if sc.StoreInterval != 500*time.Second {
		t.Errorf("Expected StoreInterval to be 500s, got %v", sc.StoreInterval)
	}
	if sc.FileStoragePath != "/tmp/storage.json" {
		t.Errorf("Expected FileStoragePath to be '/tmp/storage.json', got '%s'", sc.FileStoragePath)
	}
	if sc.Restore != false {
		t.Errorf("Expected Restore to be false, got %v", sc.Restore)
	}
	if sc.DatabaseDsn != "postgres://user:pass@localhost/db" {
		t.Errorf("Expected DatabaseDsn to be 'postgres://user:pass@localhost/db', got '%s'", sc.DatabaseDsn)
	}
	if sc.HashKey != "secret" {
		t.Errorf("Expected HashKey to be 'secret', got '%s'", sc.HashKey)
	}
	if sc.CryptoKey != "/path/to/key" {
		t.Errorf("Expected CryptoKey to be '/path/to/key', got '%s'", sc.CryptoKey)
	}
}

func TestNewServerConfig_EnvVars(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	resetFlags()

	err := os.Setenv("ADDRESS", "localhost:8082")
	if err != nil {
		t.Errorf("Failed to set environment variable ADDRESS: %v", err)
	}
	err = os.Setenv("STORE_INTERVAL", "600s")
	if err != nil {
		t.Errorf("Failed to set environment variable STORE_INTERVAL: %v", err)
	}
	err = os.Setenv("FILE_STORAGE_PATH", "/var/storage.json")
	if err != nil {
		t.Errorf("Failed to set environment variable FILE_STORAGE_PATH: %v", err)
	}
	err = os.Setenv("RESTORE", "false")
	if err != nil {
		t.Errorf("Failed to set environment variable RESTORE: %v", err)
	}
	err = os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/envdb")
	if err != nil {
		t.Errorf("Failed to set environment variable DATABASE_DSN: %v", err)
	}
	err = os.Setenv("KEY", "envsecret")
	if err != nil {
		t.Errorf("Failed to set environment variable KEY: %v", err)
	}
	err = os.Setenv("CRYPTO_KEY", "/env/path/to/key")
	if err != nil {
		t.Errorf("Failed to set environment variable CRYPTO_KEY: %v", err)
	}
	err = os.Setenv("CONFIG", "")
	if err != nil {
		t.Errorf("Failed to set environment variable CONFIG: %v", err)
	}

	sc := NewServerConfig()

	if sc.Host != "localhost:8082" {
		t.Errorf("Expected Host to be 'localhost:8082', got '%s'", sc.Host)
	}
	if sc.StoreInterval != 600*time.Second {
		t.Errorf("Expected StoreInterval to be 600s, got %v", sc.StoreInterval)
	}
	if sc.FileStoragePath != "/var/storage.json" {
		t.Errorf("Expected FileStoragePath to be '/var/storage.json', got '%s'", sc.FileStoragePath)
	}
	if sc.Restore != false {
		t.Errorf("Expected Restore to be false, got %v", sc.Restore)
	}
	if sc.DatabaseDsn != "postgres://user:pass@localhost/envdb" {
		t.Errorf("Expected DatabaseDsn to be 'postgres://user:pass@localhost/envdb', got '%s'", sc.DatabaseDsn)
	}
	if sc.HashKey != "envsecret" {
		t.Errorf("Expected HashKey to be 'envsecret', got '%s'", sc.HashKey)
	}
	if sc.CryptoKey != "/env/path/to/key" {
		t.Errorf("Expected CryptoKey to be '/env/path/to/key', got '%s'", sc.CryptoKey)
	}
	if sc.Config != "" {
		t.Errorf("Expected Config to be '', got '%s'", sc.Config)
	}
}

func resetFlags() {
	os.Clearenv()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	resetForTesting(nil)
}

// ResetForTesting clears all flag state and sets the usage function as directed.
// After calling ResetForTesting, parse errors in flag handling will not
// exit the program.
func resetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	flag.CommandLine.Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
		flag.PrintDefaults()
	}
	flag.Usage = usage
}
