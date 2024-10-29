package config

import (
	"flag"
	"os"
	"strconv"
)

// ServerConfig хранит информацию о конфгурации сервера.
type ServerConfig struct {
	// Host - адрес эндпоинта HTTP-сервера.
	Host string
	// StoreInternal - интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск.
	StoreInternal int
	// FileStoragePath - полное имя файла, куда сохраняются текущие значения.
	FileStoragePath string
	// Restore - (булево) определяет, загружать или нет ранее сохранённые значения из файла при старте сервера.
	Restore bool
	// DatabaseDsn - строка с адресом подключения к БД.
	DatabaseDsn string
	// HashKey - ключ для хеша.
	HashKey string
	// CryptoKey - путь до файла с публичным ключом.
	CryptoKey string
}

// NewServerConfig создаёт и возвращает новый экземпляр ServerConfig, инициализированный с помощью флагов.
func NewServerConfig() ServerConfig {
	sc := ServerConfig{}
	sc.parseFlag()
	return sc
}

func (sc *ServerConfig) parseFlag() {
	flag.StringVar(&sc.Host, "a", DefaultHost, "address and port to run server")
	flag.IntVar(&sc.StoreInternal, "i", 300, "store interval")
	flag.StringVar(&sc.FileStoragePath, "f", "", "file storage path")
	flag.BoolVar(&sc.Restore, "r", true, "restore")
	flag.StringVar(&sc.DatabaseDsn, "d", "", "database dsn")
	flag.StringVar(&sc.HashKey, "k", "", "hash key")
	flag.StringVar(&sc.CryptoKey, "crypto-key", DefaultCryptoKeyServer, "crypto key")
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
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		sc.DatabaseDsn = envDatabaseDsn
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		sc.HashKey = envHashKey
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		switch envRestore {
		case "true":
			sc.Restore = true
		case "false":
			sc.Restore = false
		}
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		sc.CryptoKey = envCryptoKey
	}
}
