package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/moonicy/gometrics/pkg/gzip"
)

// GzipMiddleware возвращает middleware, который обрабатывает сжатие и декомпрессию HTTP-запросов и ответов с использованием gzip.
// Он проверяет заголовки запроса и при необходимости сжимает или декомпрессирует данные, устанавливая соответствующие заголовки.
func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ow := res

		acceptEncoding := req.Header.Get("Accept-Encoding")
		contentType := req.Header.Get("Content-Type")
		accept := req.Header.Get("Accept")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			if contentType == "application/json" || contentType == "text/html" || contentType == "text/plain" || accept == "html/text" {
				res.Header().Set("Content-Encoding", "gzip")

				cw := gzip.NewCompressWriter(res)
				ow = cw
				defer func(cw *gzip.CompressWriter) {
					err := cw.Close()
					if err != nil {
						log.Fatal(err)
					}
				}(cw)
			}
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := gzip.NewCompressReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer func(cr *gzip.CompressReader) {
				err = cr.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(cr)
		}

		h.ServeHTTP(ow, req)
	})
}
