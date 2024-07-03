package config

import "strings"

const DefaultHost = "localhost:8080"

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
