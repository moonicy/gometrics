package middlewares

import (
	"github.com/moonicy/gometrics/pkg/gzip"
	"net/http"
	"strings"
)

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
				defer cw.Close()
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
			defer cr.Close()
		}

		h.ServeHTTP(ow, req)
	})
}
