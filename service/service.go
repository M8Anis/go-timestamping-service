package service

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

var signingKey *ecdsa.PrivateKey
var signingCertificate *x509.Certificate

func init() {
	pemKey, err := os.ReadFile("./certs/ts-key.pem")
	if err != nil {
		log.Fatalf("Private key cannot be read: %s", err)
	}
	key, _ := pem.Decode(pemKey)
	if key == nil {
		log.Fatalf("Private key cannot be decoded: %s", err)
	}
	signingKey, err = x509.ParseECPrivateKey(key.Bytes)
	if err != nil {
		log.Fatalf("EC private key cannot be parsed: %s", err)
	}

	pemCert, err := os.ReadFile("./certs/ts-crt.pem")
	if err != nil {
		log.Fatalf("Certificate cannot be read: %s", err)
	}
	cert, _ := pem.Decode(pemCert)
	if key == nil {
		log.Fatalf("Certificate cannot be decoded: %s", err)
	}
	signingCertificate, err = x509.ParseCertificate(cert.Bytes)
	if err != nil {
		log.Fatalf("Certificate cannot be parsed: %s", err)
	}
}
