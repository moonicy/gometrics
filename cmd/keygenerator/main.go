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
	err = pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Fatal(err)
	}

	filePrivate, err := os.Create("keys/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	_, err = filePrivate.Write(privateKeyPEM.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	defer func(filePrivate *os.File) {
		err = filePrivate.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(filePrivate)

	filePublic, err := os.Create("keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}
	_, err = filePublic.Write(publicKeyPEM.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	defer func(filePublic *os.File) {
		err = filePublic.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(filePublic)

}
