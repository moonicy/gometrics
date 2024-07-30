package config

import "strings"

const (
	DefaultHost           = "localhost:8080"
	DefaultReportInterval = 20
	DefaultPollInterval   = 2
	DefaultHashKey        = ""
	DefaultRateLimit      = 0
)

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
