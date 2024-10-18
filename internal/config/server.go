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

// ServerConfig хранит информацию о конфгурации сервера.
type ServerConfig struct {
	// Host - адрес эндпоинта HTTP-сервера.
	Host string `json:"address"`
	// StoreInterval - интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск.
	StoreInterval time.Duration `json:"store_interval"`
	// FileStoragePath - полное имя файла, куда сохраняются текущие значения.
	FileStoragePath string `json:"store_file"`
	// Restore - (булево) определяет, загружать или нет ранее сохранённые значения из файла при старте сервера.
	Restore bool `json:"restore"`
	// DatabaseDsn - строка с адресом подключения к БД.
	DatabaseDsn string `json:"database_dsn"`
	// HashKey - ключ для хеша.
	HashKey string
	// CryptoKey - путь до файла с публичным ключом.
	CryptoKey string `json:"crypto_key"`
	// Config - путь до файла конфигурации.
	Config string
}

// NewServerConfig создаёт и возвращает новый экземпляр ServerConfig, инициализированный с помощью флагов.
func NewServerConfig() ServerConfig {
	sc := ServerConfig{}
	sc.parseFlag()
	return sc
}

func (sc *ServerConfig) parseFlag() {
	var scFlags ServerConfig
	var restore string
	flag.StringVar(&scFlags.Host, "a", DefaultHost, "address and port to run server")
	flag.DurationVar(&scFlags.StoreInterval, "i", 300*time.Second, "store interval")
	flag.StringVar(&scFlags.FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&restore, "r", "", "restore")
	flag.StringVar(&scFlags.DatabaseDsn, "d", "", "database dsn")
	flag.StringVar(&scFlags.HashKey, "k", "", "hash key")
	flag.StringVar(&scFlags.CryptoKey, "crypto-key", DefaultCryptoKeyServer, "crypto key")
	flag.StringVar(&scFlags.Config, "c", "", "file config")
	flag.StringVar(&sc.Config, "config", "", "file config")
	flag.Parse()

	if scFlags.Config != "" {
		sc.Config = scFlags.Config
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		sc.Config = envConfig
	}
	if sc.Config != "" {
		file, err := os.ReadFile(sc.Config)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(file, &sc)
		if err != nil {
			log.Fatal(err)
		}
	}

	if scFlags.Host != "" {
		sc.Host = scFlags.Host
	}
	if scFlags.StoreInterval > 0 {
		sc.StoreInterval = scFlags.StoreInterval
	}
	if scFlags.FileStoragePath != "" {
		sc.FileStoragePath = scFlags.FileStoragePath
	}
	if restore != "" {
		switch restore {
		case "false":
			sc.Restore = false
		default:
			sc.Restore = true
		}
	}
	if scFlags.DatabaseDsn != "" {
		sc.DatabaseDsn = scFlags.DatabaseDsn
	}
	if scFlags.HashKey != "" {
		sc.HashKey = scFlags.HashKey
	}
	if scFlags.CryptoKey != "" {
		sc.CryptoKey = scFlags.CryptoKey
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		sc.Host = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		str := strings.Trim(envStoreInterval, "\"")
		if i, err := strconv.Atoi(str); err == nil {
			sc.StoreInterval = time.Duration(i)
		}

		if dur, err := time.ParseDuration(str); err == nil {
			sc.StoreInterval = dur
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		sc.FileStoragePath = envFileStoragePath
	}
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		sc.DatabaseDsn = envDatabaseDsn
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		sc.HashKey = envHashKey
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		switch envRestore {
		case "false":
			sc.Restore = false
		default:
			sc.Restore = true
		}
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		sc.CryptoKey = envCryptoKey
	}
}
