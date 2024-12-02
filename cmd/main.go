package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/service"
)

var chainLength int = 0
var fullCertChain []*x509.Certificate // Chain with `signingCertificate`
var certChain []*x509.Certificate     // Chain without `signingCertificate`

var signingKey crypto.Signer // RSA or EC private key
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
		ecParseErr := err
		signingKey, err = x509.ParsePKCS1PrivateKey(key.Bytes)
		if err != nil {
			log.Fatalf("Private key cannot be parsed. EC: %s; RSA: %s", ecParseErr, err)
		}
	}

	certFile := "./certs/ts-crt.pem"
	// Chain exist check
	if _, err := os.Stat("./certs/ts-crt_chain.pem"); !errors.Is(err, os.ErrNotExist) {
		certFile = "./certs/ts-crt_chain.pem"
	}
	pemCertChain, err := os.ReadFile(certFile)
	for {
		pemCert, rest := pem.Decode(pemCertChain)
		if pemCert == nil {
			log.Fatalf("Certificate cannot be decoded: %s", err)
		}
		cert, err := x509.ParseCertificate(pemCert.Bytes)
		if err != nil {
			log.Fatalf("Certificate cannot be parsed: %s", err)
		}

		fullCertChain = append(fullCertChain, cert)
		chainLength++

		if chainLength > 1 {
			curCert := fullCertChain[chainLength-1]
			prevCert := fullCertChain[chainLength-2]
			if err := prevCert.CheckSignatureFrom(curCert); err != nil {
				log.Fatalf("Invalid certificate chain: %s (`%s`, parent: `%s`)",
					err, curCert.Subject, prevCert.Subject,
				)
			}
		}

		if len(rest) == 0 {
			break
		}
		pemCertChain = rest
	}

	signingCertificate = fullCertChain[0]
	certChain = fullCertChain[1:]
}

func main() {
	service.Serve(
		chainLength, fullCertChain,
		certChain, signingCertificate, signingKey,
	)
}
