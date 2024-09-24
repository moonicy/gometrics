package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestCalcHash(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		key      string
		expected string
	}{
		{
			name:     "Test with body and key",
			body:     []byte("test body"),
			key:      "secret_key",
			expected: calculateExpectedHash("secret_key", "test body"),
		},
		{
			name:     "Test with empty body",
			body:     []byte(""),
			key:      "secret_key",
			expected: calculateExpectedHash("secret_key", ""),
		},
		{
			name:     "Test with empty key",
			body:     []byte("test body"),
			key:      "",
			expected: calculateExpectedHash("", "test body"),
		},
		{
			name:     "Test with empty body and key",
			body:     []byte(""),
			key:      "",
			expected: calculateExpectedHash("", ""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalcHash(tc.body, tc.key)
			if result != tc.expected {
				t.Errorf("CalcHash() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func calculateExpectedHash(key, body string) string {
	sha := sha256.New()
	sha.Write([]byte(key))
	sha.Write([]byte(body))
	return hex.EncodeToString(sha.Sum(nil))
}
