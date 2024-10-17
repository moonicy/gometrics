package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

// Compress сжимает переданный срез байт src и возвращает io.Reader с сжатыми данными или ошибку.
func Compress(src []byte) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	defer func(zb *gzip.Writer) {
		err := zb.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(zb)

	_, err := zb.Write(src)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
