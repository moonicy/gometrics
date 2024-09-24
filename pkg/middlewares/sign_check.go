package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"strings"

	sign "github.com/moonicy/gometrics/pkg/hash"
)

type signResponseWriter struct {
	http.ResponseWriter
	hash hash.Hash
}

func newSignResponseWriter(w http.ResponseWriter, hash hash.Hash, key string) *signResponseWriter {
	hash.Write([]byte(key))
	return &signResponseWriter{
		ResponseWriter: w,
		hash:           hash,
	}
}

func (srw *signResponseWriter) WriteHeader(code int) {
	srw.ResponseWriter.WriteHeader(code)
}

func (srw *signResponseWriter) Write(b []byte) (int, error) {
	size, err := srw.ResponseWriter.Write(b)
	srw.hash.Write(b)
	return size, err
}

func (srw *signResponseWriter) GetHashSum() string {
	return hex.EncodeToString(srw.hash.Sum(nil))
}

// SignCheckMiddleware возвращает middleware, который проверяет хэш подписи запроса и добавляет подпись к ответу.
// Он использует переданный ключ для вычисления SHA256-хэша тела запроса и сравнивает его с хэшем из заголовка "HashSHA256".
// Если хэши не совпадают, возвращает HTTP 400 Bad Request.
// В ответ добавляет заголовок "HashSHA256" с хэшем тела ответа.
func SignCheckMiddleware(key string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if key == "" {
				handler.ServeHTTP(res, req)
				return
			}
			hashHeader := req.Header.Get("HashSHA256")
			contentType := req.Header.Get("Content-Type")
			body, err := io.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			bs := sign.CalcHash(body, key)
			if hashHeader != bs && hashHeader != "" {
				if strings.Contains(contentType, "application/json") {
					res.Header().Set("Content-Type", "application/json")
				}
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			srw := newSignResponseWriter(res, sha256.New(), key)

			req.Body = io.NopCloser(bytes.NewReader(body))

			handler.ServeHTTP(srw, req)

			res.Header().Set("HashSHA256", srw.GetHashSum())
		})
	}
}
