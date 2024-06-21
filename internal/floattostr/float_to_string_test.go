package floattostr

import "testing"

func TestFloatToString(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Positive integer", 1, "1"},
		{"Negative integer", -1, "-1"},
		{"Positive float", 1.23, "1.23"},
		{"Negative float", -1.23, "-1.23"},
		{"Large number", 1234567890.1234567, "1234567890.1234567"},
		{"Small number", 0.000000123456789, "0.000000123456789"},
		{"Zero at the end", 917305.527000, "917305.527"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FloatToString(tt.input)
			if result != tt.expected {
				t.Errorf("FloatToString(%v) got %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
