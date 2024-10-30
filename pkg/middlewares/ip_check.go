package middlewares

import (
	"net"
	"net/http"
)

// IPInTrustedSubnet проверяет, принадлежит ли IP-адрес доверенной подсети
func IPInTrustedSubnet(ip, trustedSubnet string) bool {
	_, subnet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false
	}
	parsedIP := net.ParseIP(ip)
	return subnet.Contains(parsedIP)
}

func IPCheckMiddleware(trustedSubnet string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if trustedSubnet != "" {
				ipStr := req.Header.Get("X-Real-IP")
				if ipStr == "" {
					res.WriteHeader(http.StatusForbidden)
					return
				}
				if !IPInTrustedSubnet(ipStr, trustedSubnet) {
					res.WriteHeader(http.StatusForbidden)
					return
				}
				handler.ServeHTTP(res, req)
			}
		})
	}
}
