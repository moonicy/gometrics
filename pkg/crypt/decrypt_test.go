package crypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	privateKeyPEM = nil

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	tmpFile, err := os.CreateTemp("", "private_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(privateKeyPEM)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	originalData := []byte("test data")
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, originalData)
	assert.NoError(t, err)

	decryptedData, err := Decrypt(tmpFile.Name(), encryptedData)
	assert.NoError(t, err)
	assert.Equal(t, originalData, decryptedData)
}

func TestDecrypt_InvalidPrivateKeyPath(t *testing.T) {
	privateKeyPEM = nil

	_, err := Decrypt("invalid_path.pem", []byte("some data"))
	assert.Error(t, err)
}

func TestDecrypt_InvalidPEMBlock(t *testing.T) {
	privateKeyPEM = nil

	tmpFile, err := os.CreateTemp("", "invalid_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("INVALID PEM DATA"))
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	_, err = Decrypt(tmpFile.Name(), []byte("some data"))
	assert.Error(t, err)
	assert.Equal(t, "failed to decode PEM block containing private key", err.Error())
}

func TestDecrypt_IncorrectEncryptedData(t *testing.T) {
	privateKeyPEM = nil

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	tmpFile, err := os.CreateTemp("", "private_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(privateKeyPEM)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	_, err = Decrypt(tmpFile.Name(), []byte("invalid encrypted data"))
	assert.Error(t, err)
}

func TestDecrypt_SegmentSize(t *testing.T) {
	privateKeyPEM = nil

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	tmpFile, err := os.CreateTemp("", "private_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(privateKeyPEM)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	originalData := bytes.Repeat([]byte("a"), privateKey.Size()-11*2)
	var encryptedData []byte
	segmentSize := privateKey.Size() - 11

	for i := 0; i < len(originalData); i += segmentSize {
		end := i + segmentSize
		if end > len(originalData) {
			end = len(originalData)
		}
		segment, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, originalData[i:end])
		assert.NoError(t, err)
		encryptedData = append(encryptedData, segment...)
	}

	decryptedData, err := Decrypt(tmpFile.Name(), encryptedData)
	assert.NoError(t, err)
	assert.Equal(t, originalData, decryptedData)
}
