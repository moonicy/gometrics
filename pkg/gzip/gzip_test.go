package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Compress non-empty data",
			input:   []byte("Hello, world!"),
			wantErr: false,
		},
		{
			name:    "Compress empty data",
			input:   []byte(""),
			wantErr: false,
		},
		{
			name:    "Compress large data",
			input:   bytes.Repeat([]byte("a"), 10000),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			compressedReader, err := Compress(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			gzipReader, err := gzip.NewReader(compressedReader)
			if err != nil {
				t.Errorf("Failed to create gzip reader: %v", err)
				return
			}
			defer gzipReader.Close()

			decompressedData, err := io.ReadAll(gzipReader)
			if err != nil {
				t.Errorf("Failed to read decompressed data: %v", err)
				return
			}

			assert.Equal(t, tc.input, decompressedData)
		})
	}
}
