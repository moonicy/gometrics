package main

import (
	"flag"
	"strings"
)

var flagRunAddr string
var flagReportInterval string
var flagPollInterval string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.StringVar(&flagReportInterval, "r", "10", "report interval")
	flag.StringVar(&flagPollInterval, "p", "2", "poll interval")
	flag.Parse()
}

func parseURI(uri string) string {
	str := strings.Split(uri, ":")
	if len(str) == 1 {
		return "http://localhost" + uri
	}
	if len(str) < 3 {
		return "http://" + uri
	}
	return uri
}
