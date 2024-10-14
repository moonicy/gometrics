package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	filePrivate, err := os.Create("keys/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	filePrivate.Write(privateKeyPEM.Bytes())

	defer filePrivate.Close()

	filePublic, err := os.Create("keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}
	filePublic.Write(publicKeyPEM.Bytes())

	defer filePublic.Close()

}
