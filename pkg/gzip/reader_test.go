package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockReadCloser struct {
	io.Reader
	closed bool
}

func (mrc *mockReadCloser) Close() error {
	mrc.closed = true
	return nil
}

func TestCompressReader_Close(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write([]byte("Test data"))
	if err != nil {
		t.Fatalf("Failed to write gzip data: %v", err)
	}
	err = gw.Close()
	if err != nil {
		log.Fatal(err)
	}

	mockRC := &mockReadCloser{Reader: &buf}

	cr, err := NewCompressReader(mockRC)
	if err != nil {
		t.Fatalf("Failed to create CompressReader: %v", err)
	}

	err = cr.Close()
	if err != nil {
		t.Fatalf("Failed to close CompressReader: %v", err)
	}

	assert.True(t, mockRC.closed)
}

func TestCompressReader_Read(t *testing.T) {
	originalData := []byte("Hello, World!")

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(originalData)
	if err != nil {
		t.Fatalf("Failed to write gzip data: %v", err)
	}
	err = gw.Close()
	if err != nil {
		log.Fatal(err)
	}

	cr, err := NewCompressReader(io.NopCloser(&buf))
	if err != nil {
		t.Fatalf("Failed to create CompressReader: %v", err)
	}
	defer func(cr *CompressReader) {
		err = cr.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(cr)

	decompressedData, err := io.ReadAll(cr)
	if err != nil {
		t.Fatalf("Failed to read decompressed data: %v", err)
	}

	assert.Equal(t, originalData, decompressedData)
}

func TestNewCompressReader_Error(t *testing.T) {
	badData := bytes.NewBuffer([]byte("invalid gzip data"))
	rc := io.NopCloser(badData)

	cr, err := NewCompressReader(rc)
	assert.Nil(t, cr)
	assert.Error(t, err)
}
