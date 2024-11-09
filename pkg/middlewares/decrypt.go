package middlewares

import (
	"bytes"
	"github.com/moonicy/gometrics/pkg/crypt"
	"io"
	"log"
	"net/http"
)

// CryptMiddleware .
func CryptMiddleware(publicKeyPath, privateKeyPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Failed to read request body", http.StatusInternalServerError)
					return
				}

				if len(bodyBytes) != 0 {
					decryptedBody, err := crypt.Decrypt(privateKeyPath, bodyBytes)
					if err != nil {
						http.Error(w, "Failed to decrypt request", http.StatusInternalServerError)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(decryptedBody))
				}
			}

			responseRecorder := &ResponseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
			next.ServeHTTP(responseRecorder, r)

			encryptedResponse, err := crypt.Encrypt(publicKeyPath, responseRecorder.body.Bytes())
			if err != nil {
				http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/octet-stream")
			_, err = w.Write(encryptedResponse)
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}

// ResponseRecorder для перехвата данных ответа
type ResponseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
