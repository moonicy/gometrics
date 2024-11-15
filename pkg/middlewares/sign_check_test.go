package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	mhash "github.com/moonicy/gometrics/pkg/hash"
)

func TestSignCheckMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Test response"))
		if err != nil {
			log.Fatal(err)
		}
	})

	tests := []struct {
		name               string
		key                string
		requestBody        []byte
		requestHashHeader  string
		contentType        string
		expectedStatusCode int
		expectHashHeader   bool
	}{
		{
			name:               "Empty key",
			key:                "",
			requestBody:        []byte("Test request body"),
			requestHashHeader:  "",
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectHashHeader:   false,
		},
		{
			name:               "Valid hash header",
			key:                "secret",
			requestBody:        []byte("Test request body"),
			requestHashHeader:  mhash.CalcHash([]byte("Test request body"), "secret"),
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectHashHeader:   true,
		},
		{
			name:               "Invalid hash header",
			key:                "secret",
			requestBody:        []byte("Test request body"),
			requestHashHeader:  "invalidhash",
			contentType:        "application/json",
			expectedStatusCode: http.StatusBadRequest,
			expectHashHeader:   false,
		},
		{
			name:               "Missing hash header",
			key:                "secret",
			requestBody:        []byte("Test request body"),
			requestHashHeader:  "",
			contentType:        "application/json",
			expectedStatusCode: http.StatusOK,
			expectHashHeader:   true,
		},
		{
			name:               "Non-JSON content type",
			key:                "secret",
			requestBody:        []byte("Test request body"),
			requestHashHeader:  "invalidhash",
			contentType:        "text/plain",
			expectedStatusCode: http.StatusBadRequest,
			expectHashHeader:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", tc.contentType)
			if tc.requestHashHeader != "" {
				req.Header.Set("HashSHA256", tc.requestHashHeader)
			}

			rr := httptest.NewRecorder()

			middleware := SignCheckMiddleware(tc.key)(testHandler)

			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatusCode, rr.Code)

			if tc.expectHashHeader && rr.Code == http.StatusOK {
				h := sha256.New()
				h.Write([]byte(tc.key))
				h.Write([]byte("Test response"))
				expectedResponseHash := hex.EncodeToString(h.Sum(nil))

				responseHashHeader := rr.Header().Get("HashSHA256")

				assert.Equal(t, expectedResponseHash, responseHashHeader)
			} else {
				responseHashHeader := rr.Header().Get("HashSHA256")
				assert.Empty(t, responseHashHeader)
			}
		})
	}
}

func TestSignResponseWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	var hasher = sha256.New()

	srw := &signResponseWriter{
		ResponseWriter: recorder,
		hash:           hasher,
	}

	srw.WriteHeader(http.StatusAccepted)

	if recorder.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, recorder.Code)
	}
}
