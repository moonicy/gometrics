package config

import "testing"

func TestParseURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No colon in URI",
			input:    "localhost",
			expected: "http://localhostlocalhost",
		},
		{
			name:     "One colon in URI",
			input:    "localhost:8080",
			expected: "http://localhost:8080",
		},
		{
			name:     "Two colons in URI",
			input:    "http://localhost:8080",
			expected: "http://localhost:8080",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "http://localhost",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseURI(tc.input)
			if result != tc.expected {
				t.Errorf("ParseURI(%q) = %q; want %q", tc.input, result, tc.expected)
			}
		})
	}
}
