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

const MaxSegmentSize = 214 // Максимальный размер сегмента для шифрования

var publicKeyPEM []byte

func Encrypt(publicKeyPath string, data []byte) ([]byte, error) {
	if publicKeyPEM == nil {
		pkPEM, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, err
		}
		publicKeyPEM = pkPEM
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	var encryptedBuffer bytes.Buffer

	for i := 0; i < len(data); i += MaxSegmentSize {
		end := i + MaxSegmentSize
		if end > len(data) {
			end = len(data)
		}
		segment := data[i:end]

		encryptedSegment, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, segment)
		if err != nil {
			return nil, err
		}

		encryptedBuffer.Write(encryptedSegment)
	}

	return encryptedBuffer.Bytes(), nil
}
