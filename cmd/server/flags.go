package main

import (
	"flag"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/http"
	"os"
	"strconv"
)

func parseFlag() config.ServerConfig {
	cfg := config.ServerConfig{}

	flag.StringVar(&cfg.Host, "a", http.DefaultHost, "address and port to run server")
	flag.IntVar(&cfg.StoreInternal, "i", 300, "store interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Host = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		cfg.StoreInternal, _ = strconv.Atoi(envStoreInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		switch envRestore {
		case "true":
			cfg.Restore = true
		case "false":
			cfg.Restore = false
		}
	}
	return cfg
}
