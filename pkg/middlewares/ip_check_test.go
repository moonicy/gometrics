package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPInTrustedSubnet(t *testing.T) {
	tests := []struct {
		name          string
		ip            string
		trustedSubnet string
		want          bool
	}{
		{
			name:          "IP in subnet",
			ip:            "192.168.1.5",
			trustedSubnet: "192.168.1.0/24",
			want:          true,
		},
		{
			name:          "IP not in subnet",
			ip:            "192.168.2.5",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "Invalid IP",
			ip:            "invalid_ip",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "Invalid subnet",
			ip:            "192.168.1.5",
			trustedSubnet: "invalid_subnet",
			want:          false,
		},
		{
			name:          "IPv6 in subnet",
			ip:            "2001:db8::1",
			trustedSubnet: "2001:db8::/32",
			want:          true,
		},
		{
			name:          "IPv6 not in subnet",
			ip:            "2001:db9::1",
			trustedSubnet: "2001:db8::/32",
			want:          false,
		},
		{
			name:          "Empty IP and subnet",
			ip:            "",
			trustedSubnet: "",
			want:          false,
		},
		{
			name:          "Empty IP",
			ip:            "",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "Empty subnet",
			ip:            "192.168.1.5",
			trustedSubnet: "",
			want:          false,
		},
		{
			name:          "IPv4-mapped IPv6 in subnet",
			ip:            "::ffff:192.168.1.5",
			trustedSubnet: "192.168.1.0/24",
			want:          true,
		},
		{
			name:          "IPv4-mapped IPv6 not in subnet",
			ip:            "::ffff:192.168.2.5",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "Invalid CIDR notation",
			ip:            "192.168.1.5",
			trustedSubnet: "192.168.1.0/33",
			want:          false,
		},
		{
			name:          "Valid IP, invalid subnet",
			ip:            "192.168.1.5",
			trustedSubnet: "300.300.300.0/24",
			want:          false,
		},
		{
			name:          "Invalid IP, valid subnet",
			ip:            "999.999.999.999",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "IP at subnet boundary",
			ip:            "192.168.1.0",
			trustedSubnet: "192.168.1.0/24",
			want:          true,
		},
		{
			name:          "IP at subnet broadcast",
			ip:            "192.168.1.255",
			trustedSubnet: "192.168.1.0/24",
			want:          true,
		},
		{
			name:          "IP outside subnet range",
			ip:            "192.168.2.0",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "IPv6 zero address",
			ip:            "::",
			trustedSubnet: "::/0",
			want:          true,
		},
		{
			name:          "IPv4 zero address",
			ip:            "0.0.0.0",
			trustedSubnet: "0.0.0.0/0",
			want:          true,
		},
		{
			name:          "IPv4 in IPv6 subnet",
			ip:            "192.168.1.5",
			trustedSubnet: "2001:db8::/32",
			want:          false,
		},
		{
			name:          "IPv6 in IPv4 subnet",
			ip:            "2001:db8::1",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "Different address families",
			ip:            "::ffff:c0a8:0105", // ::ffff:192.168.1.5
			trustedSubnet: "192.168.1.0/24",
			want:          true,
		},
		{
			name:          "Different address families not in subnet",
			ip:            "::ffff:c0a8:0205", // ::ffff:192.168.2.5
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "IPv6 subnet with prefix length 128",
			ip:            "2001:db8::1",
			trustedSubnet: "2001:db8::1/128",
			want:          true,
		},
		{
			name:          "IPv6 subnet with prefix length 128 not matching",
			ip:            "2001:db8::2",
			trustedSubnet: "2001:db8::1/128",
			want:          false,
		},
		{
			name:          "IPv4 subnet with prefix length 32",
			ip:            "192.168.1.5",
			trustedSubnet: "192.168.1.5/32",
			want:          true,
		},
		{
			name:          "IPv4 subnet with prefix length 32 not matching",
			ip:            "192.168.1.6",
			trustedSubnet: "192.168.1.5/32",
			want:          false,
		},
		{
			name:          "IPv4 address in default subnet",
			ip:            "10.1.2.3",
			trustedSubnet: "0.0.0.0/0",
			want:          true,
		},
		{
			name:          "IPv6 address in default subnet",
			ip:            "2001:db8::1",
			trustedSubnet: "::/0",
			want:          true,
		},
		{
			name:          "IPv4 private address in public subnet",
			ip:            "192.168.1.5",
			trustedSubnet: "203.0.113.0/24",
			want:          false,
		},
		{
			name:          "IPv4 public address in public subnet",
			ip:            "203.0.113.5",
			trustedSubnet: "203.0.113.0/24",
			want:          true,
		},
		{
			name:          "IPv6 address not matching subnet",
			ip:            "2001:db8:1::1",
			trustedSubnet: "2001:db8::/48",
			want:          false,
		},
		{
			name:          "IPv6 address matching subnet",
			ip:            "2001:db8:0:1::1",
			trustedSubnet: "2001:db8::/48",
			want:          true,
		},
		{
			name:          "Invalid IP and invalid subnet",
			ip:            "invalid_ip",
			trustedSubnet: "invalid_subnet",
			want:          false,
		},
		{
			name:          "IP with extra spaces",
			ip:            " 192.168.1.5 ",
			trustedSubnet: "192.168.1.0/24",
			want:          false,
		},
		{
			name:          "Subnet with extra spaces",
			ip:            "192.168.1.5",
			trustedSubnet: " 192.168.1.0/24 ",
			want:          false,
		},
		{
			name:          "IP in IPv6 subnet using IPv4 address",
			ip:            "192.168.1.5",
			trustedSubnet: "::ffff:192.168.1.0/112",
			want:          true,
		},
		{
			name:          "Loopback address in subnet",
			ip:            "127.0.0.1",
			trustedSubnet: "127.0.0.0/8",
			want:          true,
		},
		{
			name:          "Loopback address not in subnet",
			ip:            "127.0.0.1",
			trustedSubnet: "10.0.0.0/8",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IPInTrustedSubnet(tt.ip, tt.trustedSubnet)
			if got != tt.want {
				t.Errorf("IPInTrustedSubnet(%q, %q) = %v; want %v", tt.ip, tt.trustedSubnet, got, tt.want)
			}
		})
	}
}

func TestIPCheckMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		xRealIP       string
		wantStatus    int
	}{
		{
			name:          "Allowed IP in trusted subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Forbidden IP not in trusted subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.2.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Missing X-Real-IP header",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Invalid trusted subnet",
			trustedSubnet: "invalid_subnet",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "IPv6 allowed in trusted subnet",
			trustedSubnet: "2001:db8::/32",
			xRealIP:       "2001:db8::1",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "IPv6 forbidden not in trusted subnet",
			trustedSubnet: "2001:db8::/32",
			xRealIP:       "2001:db9::1",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "IPv4-mapped IPv6 in subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "::ffff:192.168.1.5",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "IPv4-mapped IPv6 not in subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "::ffff:192.168.2.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Invalid X-Real-IP header",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "invalid_ip",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Handler not called on forbidden",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.2.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Handler called on allowed",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.1.10",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Trusted subnet with extra spaces",
			trustedSubnet: " 192.168.1.0/24 ",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "X-Real-IP with extra spaces",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       " 192.168.1.5 ",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "IPv6 address in IPv4 subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "2001:db8::1",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "IPv4 address in IPv6 subnet",
			trustedSubnet: "2001:db8::/32",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "No handler called on missing X-Real-IP",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Multiple X-Real-IP headers, first is valid",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.1.5,192.168.2.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Multiple X-Real-IP headers, first is invalid",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.2.5,192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "X-Real-IP header with port",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.1.5:8080",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "X-Real-IP header with CIDR notation",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "192.168.1.5/24",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Forbidden when trusted subnet invalid",
			trustedSubnet: "invalid_subnet",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Trusted subnet is 0.0.0.0/0 (all IPv4 addresses)",
			trustedSubnet: "0.0.0.0/0",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Trusted subnet is ::/0 (all IPv6 addresses)",
			trustedSubnet: "::/0",
			xRealIP:       "2001:db8::1",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Trusted subnet is ::/0, IPv4 address",
			trustedSubnet: "::/0",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Trusted subnet is 0.0.0.0/0, IPv6 address",
			trustedSubnet: "0.0.0.0/0",
			xRealIP:       "2001:db8::1",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Valid IP, invalid subnet",
			trustedSubnet: "300.300.300.0/24",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Invalid IP, valid subnet",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "999.999.999.999",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Trusted subnet with prefix length 32",
			trustedSubnet: "192.168.1.5/32",
			xRealIP:       "192.168.1.5",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Trusted subnet with prefix length 32 not matching",
			trustedSubnet: "192.168.1.5/32",
			xRealIP:       "192.168.1.6",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "IPv6 trusted subnet with prefix length 128",
			trustedSubnet: "2001:db8::1/128",
			xRealIP:       "2001:db8::1",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "IPv6 trusted subnet with prefix length 128 not matching",
			trustedSubnet: "2001:db8::1/128",
			xRealIP:       "2001:db8::2",
			wantStatus:    http.StatusForbidden,
		},
		{
			name:          "Loopback address in trusted subnet",
			trustedSubnet: "127.0.0.0/8",
			xRealIP:       "127.0.0.1",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Loopback address not in trusted subnet",
			trustedSubnet: "10.0.0.0/8",
			xRealIP:       "127.0.0.1",
			wantStatus:    http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false

			// Test handler that sets called to true
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			// Wrap the handler with middleware
			middleware := IPCheckMiddleware(tt.trustedSubnet)
			handlerToTest := middleware(nextHandler)

			// Create request
			req := httptest.NewRequest("GET", "/", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			// Record the response
			rr := httptest.NewRecorder()

			// Serve the request
			handlerToTest.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.wantStatus {
				t.Errorf("Status code = %v; want %v", rr.Code, tt.wantStatus)
			}

			// Check if handler was called
			if tt.wantStatus == http.StatusOK && !called {
				t.Errorf("Handler was not called but should have been")
			} else if tt.wantStatus != http.StatusOK && called {
				t.Errorf("Handler was called but should not have been")
			}
		})
	}
}
