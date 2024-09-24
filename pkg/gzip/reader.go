package gzip

import (
	"compress/gzip"
	"io"
)

// CompressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера декомпрессировать получаемые от клиента данные.
type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader создаёт новый CompressReader для чтения и декомпрессии данных из io.ReadCloser.
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read читает и декомпрессирует данные из внутреннего gzip.Reader.
func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает внутренние io.ReadCloser и gzip.Reader.
func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
