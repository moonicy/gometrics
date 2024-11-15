package crypt

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	publicKeyPEM = nil

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	tmpFile, err := os.CreateTemp("", "public_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(publicKeyPEM)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	data := []byte("test data for encryption")

	encryptedData, err := Encrypt(tmpFile.Name(), data)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData)

	var decryptedBuffer bytes.Buffer
	segmentSize := privateKey.Size()

	for i := 0; i < len(encryptedData); i += segmentSize {
		end := i + segmentSize
		if end > len(encryptedData) {
			end = len(encryptedData)
		}
		encryptedSegment := encryptedData[i:end]

		decryptedSegment, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedSegment)
		assert.NoError(t, err)
		decryptedBuffer.Write(decryptedSegment)
	}

	assert.Equal(t, data, decryptedBuffer.Bytes())
}

func TestEncrypt_InvalidPublicKeyPath(t *testing.T) {
	publicKeyPEM = nil

	_, err := Encrypt("invalid_path.pem", []byte("some data"))
	assert.Error(t, err)
}

func TestEncrypt_InvalidPEMBlock(t *testing.T) {
	publicKeyPEM = nil

	tmpFile, err := os.CreateTemp("", "invalid_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("INVALID PEM DATA"))
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	_, err = Encrypt(tmpFile.Name(), []byte("some data"))
	assert.Error(t, err)
	assert.Equal(t, "failed to decode PEM block containing public key", err.Error())
}

func TestEncrypt_NotRSAPublicKey(t *testing.T) {
	publicKeyPEM = nil

	ecPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	ecPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&ecPrivateKey.PublicKey)
	assert.NoError(t, err)

	ecPublicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ecPublicKeyBytes,
	})

	tmpFile, err := os.CreateTemp("", "non_rsa_public_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(ecPublicKeyPEM)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	_, err = Encrypt(tmpFile.Name(), []byte("some data"))
	assert.Error(t, err)
	assert.Equal(t, "not RSA public key", err.Error())
}

func TestEncrypt_LargeData(t *testing.T) {
	publicKeyPEM = nil

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	tmpFile, err := os.CreateTemp("", "public_key*.pem")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(publicKeyPEM)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	originalData := bytes.Repeat([]byte("a"), MaxSegmentSize*2)

	encryptedData, err := Encrypt(tmpFile.Name(), originalData)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData)

	var decryptedBuffer bytes.Buffer
	segmentSize := privateKey.Size()

	for i := 0; i < len(encryptedData); i += segmentSize {
		end := i + segmentSize
		if end > len(encryptedData) {
			end = len(encryptedData)
		}
		encryptedSegment := encryptedData[i:end]

		decryptedSegment, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedSegment)
		assert.NoError(t, err)
		decryptedBuffer.Write(decryptedSegment)
	}

	assert.Equal(t, originalData, decryptedBuffer.Bytes())
}
