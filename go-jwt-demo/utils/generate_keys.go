package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateFile.Close()

	err = pem.Encode(privateFile, &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer publicFile.Close()

	err = pem.Encode(publicFile, &pem.Block{
		Type: "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RSA key pair generated successfully")
}