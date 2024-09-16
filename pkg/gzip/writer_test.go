package gzip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressWriter_Header(t *testing.T) {
	tests := []struct {
		name          string
		headerKey     string
		headerValue   string
		expectedKey   string
		expectedValue string
	}{
		{
			name:          "Set Content-Type header",
			headerKey:     "Content-Type",
			headerValue:   "text/plain",
			expectedKey:   "Content-Type",
			expectedValue: "text/plain",
		},
		{
			name:          "Set Custom header",
			headerKey:     "X-Custom-Header",
			headerValue:   "CustomValue",
			expectedKey:   "X-Custom-Header",
			expectedValue: "CustomValue",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			cw := NewCompressWriter(responseRecorder)

			header := cw.Header()
			header.Set(tc.headerKey, tc.headerValue)

			cw.WriteHeader(http.StatusOK)
			cw.Write([]byte("Test"))
			cw.Close()

			result := responseRecorder.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.expectedValue, result.Header.Get(tc.expectedKey))
		})
	}
}

func TestCompressWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedHeader string
	}{
		{
			name:           "Status code 200",
			statusCode:     http.StatusOK,
			expectedHeader: "gzip",
		},
		{
			name:           "Status code 201",
			statusCode:     http.StatusCreated,
			expectedHeader: "gzip",
		},
		{
			name:           "Status code 400",
			statusCode:     http.StatusBadRequest,
			expectedHeader: "",
		},
		{
			name:           "Status code 500",
			statusCode:     http.StatusInternalServerError,
			expectedHeader: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			cw := NewCompressWriter(responseRecorder)
			defer cw.Close()

			cw.WriteHeader(tc.statusCode)
			cw.Write([]byte("Test"))

			err := cw.Close()
			assert.NoError(t, err)

			result := responseRecorder.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.expectedHeader, result.Header.Get("Content-Encoding"))
			assert.Equal(t, tc.statusCode, result.StatusCode)
		})
	}
}
