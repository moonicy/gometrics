package crypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var privateKeyPEM []byte

func Decrypt(privateKeyPath string, encryptedData []byte) ([]byte, error) {
	if privateKeyPEM == nil {
		pkPEM, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}
		privateKeyPEM = pkPEM
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var decryptedBuffer bytes.Buffer
	segmentSize := privateKey.Size()

	for i := 0; i < len(encryptedData); i += segmentSize {
		end := i + segmentSize
		if end > len(encryptedData) {
			end = len(encryptedData)
		}
		encryptedSegment := encryptedData[i:end]

		decryptedSegment, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedSegment)
		if err != nil {
			return nil, err
		}

		decryptedBuffer.Write(decryptedSegment)
	}

	return decryptedBuffer.Bytes(), nil
}
