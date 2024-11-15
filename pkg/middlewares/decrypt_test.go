package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestResponseRecorder_Write(t *testing.T) {
	bodyBuffer := &bytes.Buffer{}

	recorder := &ResponseRecorder{
		ResponseWriter: httptest.NewRecorder(),
		body:           bodyBuffer,
	}

	data := []byte("Test response data")

	n, err := recorder.Write(data)
	if err != nil {
		t.Fatalf("Write returned an error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d bytes", len(data), n)
	}

	if recorder.body.String() != string(data) {
		t.Errorf("Expected body buffer to contain %q, got %q", string(data), recorder.body.String())
	}

	response := recorder.ResponseWriter.(*httptest.ResponseRecorder)
	if response.Body.String() != string(data) {
		t.Errorf("Expected ResponseWriter to contain %q, got %q", string(data), response.Body.String())
	}
}

func TestResponseRecorder_WriteHeader(t *testing.T) {
	recorder := &ResponseRecorder{
		ResponseWriter: httptest.NewRecorder(),
	}

	recorder.WriteHeader(http.StatusAccepted)

	response := recorder.ResponseWriter.(*httptest.ResponseRecorder)
	if response.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, response.Code)
	}
}

func TestResponseRecorder_InterfaceCompliance(t *testing.T) {
	var _ http.ResponseWriter = &ResponseRecorder{}
}

func generateRSAKeys() (publicKeyPath, privateKeyPath string, cleanup func(), err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", nil, err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	privateKeyFile, err := os.CreateTemp("", "private_key*.pem")
	if err != nil {
		return "", "", nil, err
	}
	defer privateKeyFile.Close()

	_, err = privateKeyFile.Write(privateKeyPEM)
	if err != nil {
		return "", "", nil, err
	}

	publicKeyFile, err := os.CreateTemp("", "public_key*.pem")
	if err != nil {
		return "", "", nil, err
	}
	defer publicKeyFile.Close()

	_, err = publicKeyFile.Write(publicKeyPEM)
	if err != nil {
		return "", "", nil, err
	}

	cleanup = func() {
		os.Remove(privateKeyFile.Name())
		os.Remove(publicKeyFile.Name())
	}

	return publicKeyFile.Name(), privateKeyFile.Name(), cleanup, nil
}

func TestCryptMiddleware_DecryptionError(t *testing.T) {
	publicKeyPath, privateKeyPath, cleanup, err := generateRSAKeys()
	assert.NoError(t, err)
	defer cleanup()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	middleware := CryptMiddleware(publicKeyPath, privateKeyPath)
	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL, bytes.NewReader([]byte("invalid encrypted data")))
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Failed to decrypt request")
}

func TestCryptMiddleware_ReadBodyError(t *testing.T) {
	publicKeyPath, privateKeyPath, cleanup, err := generateRSAKeys()
	assert.NoError(t, err)
	defer cleanup()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	middleware := CryptMiddleware(publicKeyPath, privateKeyPath)
	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	brokenReader := &BrokenReader{}

	req, err := http.NewRequest("POST", server.URL, brokenReader)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Error(t, err)
	if resp != nil {
		defer resp.Body.Close()
	}
}

type BrokenReader struct{}

func (br *BrokenReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
