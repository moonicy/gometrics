package config

type ServerConfig struct {
	Host            string
	StoreInternal   int
	FileStoragePath string
	Restore         bool
}
