package service

import (
	"crypto/ecdsa"
	"crypto/x509"
	"log"
	"os"
)

var signingKey *ecdsa.PrivateKey
var signingCertificate *x509.Certificate

func init() {
	key, err := os.ReadFile("./testcert/ts-key.der")
	if err != nil {
		log.Fatalf("EC private key cannot be read: %s", err)
	}
	signingKey, err = x509.ParseECPrivateKey(key)
	if err != nil {
		log.Fatalf("EC private key cannot be parsed: %s", err)
	}

	cert, err := os.ReadFile("./testcert/ts-crt.der")
	if err != nil {
		log.Fatalf("Certificate cannot be read: %s", err)
	}
	signingCertificate, err = x509.ParseCertificate(cert)
	if err != nil {
		log.Fatalf("Certificate cannot be parsed: %s", err)
	}
}
