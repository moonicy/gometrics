package gzip

import (
	"compress/gzip"
	"net/http"
)

// CompressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера сжимать передаваемые данные и выставлять правильные HTTP-заголовки.
type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter создаёт и возвращает новый CompressWriter, обёртку над http.ResponseWriter для сжатия выходящих данных.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header возвращает заголовки HTTP-ответа.
func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

// Write записывает данные в сжатом формате в поток ответа.
func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader устанавливает код статуса HTTP-ответа и добавляет заголовок Content-Encoding при необходимости.
func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *CompressWriter) Close() error {
	return c.zw.Close()
}
