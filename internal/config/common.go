package config

import "strings"

// Конфигурация сервера по умолчанию.
const (
	DefaultHost           = "localhost:8080"
	DefaultReportInterval = 20
	DefaultPollInterval   = 2
	DefaultHashKey        = ""
	DefaultRateLimit      = 0
)

// ParseURI возвращает полный URI, добавляя протокол и хост по умолчанию при необходимости.
func ParseURI(uri string) string {
	str := strings.Split(uri, ":")
	if len(str) == 1 {
		return "http://localhost" + uri
	}
	if len(str) < 3 {
		return "http://" + uri
	}
	return uri
}
