package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	tests := []struct {
		name            string
		acceptEncoding  string
		contentEncoding string
		contentType     string
		expectedBody    string
		expectedHeader  string
		isBodyGzipped   bool
	}{
		{
			name:           "No gzip support",
			acceptEncoding: "",
			expectedBody:   "Hello, world!",
			expectedHeader: "",
			isBodyGzipped:  false,
		},
		{
			name:           "Gzip support with HTML",
			acceptEncoding: "gzip",
			contentType:    "text/html",
			expectedBody:   "Hello, world!",
			expectedHeader: "gzip",
			isBodyGzipped:  true,
		},
		{
			name:           "Gzip support with JSON",
			acceptEncoding: "gzip",
			contentType:    "application/json",
			expectedBody:   "Hello, world!",
			expectedHeader: "gzip",
			isBodyGzipped:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			req.Header.Set("Accept-Encoding", tc.acceptEncoding)
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			if tc.contentEncoding != "" {
				var body bytes.Buffer
				gw := gzip.NewWriter(&body)
				_, _ = gw.Write([]byte(tc.expectedBody))
				gw.Close()
				req.Body = io.NopCloser(&body)
				req.Header.Set("Content-Encoding", tc.contentEncoding)
			}

			w := httptest.NewRecorder()

			middleware := GzipMiddleware(handler)
			middleware.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if tc.expectedHeader != "" {
				assert.Equal(t, tc.expectedHeader, res.Header.Get("Content-Encoding"))
			} else {
				assert.Empty(t, res.Header.Get("Content-Encoding"))
			}

			var responseBody []byte
			if tc.isBodyGzipped {
				gr, err := gzip.NewReader(res.Body)
				assert.NoError(t, err)
				responseBody, err = io.ReadAll(gr)
				assert.NoError(t, err)
			} else {
				responseBody, _ = io.ReadAll(res.Body)
			}

			assert.Equal(t, tc.expectedBody, string(responseBody))
		})
	}
}
