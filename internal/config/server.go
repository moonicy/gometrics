package config

import (
	"flag"
	"os"
	"strconv"
)

type ServerConfig struct {
	Host            string
	StoreInternal   int
	FileStoragePath string
	Restore         bool
}

func NewServerConfig() *ServerConfig {
	sc := &ServerConfig{}
	sc.parseFlag()
	return sc
}

func (sc *ServerConfig) parseFlag() {
	flag.StringVar(&sc.Host, "a", DefaultHost, "address and port to run server")
	flag.IntVar(&sc.StoreInternal, "i", 300, "store interval")
	flag.StringVar(&sc.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&sc.Restore, "r", true, "restore")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		sc.Host = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		sc.StoreInternal, _ = strconv.Atoi(envStoreInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		sc.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		switch envRestore {
		case "true":
			sc.Restore = true
		case "false":
			sc.Restore = false
		}
	}
}
