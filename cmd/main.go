package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"slices"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/service"
	flag "github.com/spf13/pflag"
)

var (
	_ca_chain_file         string
	_cert_file, _priv_file string

	_address string
	_port    uint16
)

var (
	_ca_chain []*x509.Certificate
	_cert     *x509.Certificate
	_priv     crypto.Signer
)

func init() {
	flag.StringVarP(&_ca_chain_file, "authorities-chain", "a", "./certs/ca.pem", "Path to CAs chain of timestamper certificate")
	flag.StringVarP(&_cert_file, "certificate", "c", "./certs/cert.pem", "Path to timestamper certificate")
	flag.StringVarP(&_priv_file, "key", "k", "./private/key.pem", "Path to timestamper certificate private key")

	flag.StringVarP(&_address, "ip", "i", "127.105.35.186", "IP to run on")
	flag.Uint16VarP(&_port, "port", "p", 13916, "Port to run on")

	flag.Parse()
}

func init() {
	if certFile, err := os.ReadFile(_cert_file); err == nil {
		pemCert, _ := pem.Decode(certFile)

		if _cert, err = x509.ParseCertificate(pemCert.Bytes); err != nil {
			log.Fatalf("Cannot parse timestamper certificate: %s", err)
		}
	} else {
		log.Fatalf("Cannot read timestamper certificate: %s", err)
	}

	if (_cert.KeyUsage & x509.KeyUsageDigitalSignature) == 0 {
		log.Fatalf("Timestamper certificate can't be used for signing (No `Digital Signature` in key usage)")
	}

	if !slices.Contains(_cert.ExtKeyUsage, x509.ExtKeyUsageTimeStamping) {
		log.Fatalf("Timestamper certificate can't be used for timestamping (No `Timestamping` in extended key usage)")
	}

	if privFile, err := os.ReadFile(_priv_file); err == nil {
		pemPriv, _ := pem.Decode(privFile)

		var ecParseError, rsaParseError error

		if _priv, ecParseError = x509.ParseECPrivateKey(pemPriv.Bytes); ecParseError != nil {
			if _priv, rsaParseError = x509.ParsePKCS1PrivateKey(pemPriv.Bytes); rsaParseError != nil {
				log.Fatalf("Cannot parse timestamper private key.\nEC: %s; RSA: %s", ecParseError, rsaParseError)
			}
		}
	} else {
		log.Fatalf("Cannot read timestamper private key: %s", err)
	}

	pubKeyMatch := true
	switch _cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		pubKeyMatch = _cert.PublicKey.(*ecdsa.PublicKey).Equal(_priv.Public().(*ecdsa.PublicKey))
	case *rsa.PublicKey:
		pubKeyMatch = _cert.PublicKey.(*rsa.PublicKey).Equal(_priv.Public().(*rsa.PublicKey))
	default:
		log.Print("Certificate public key and public key in private key not checked. Unknown public key type")
	}
	if !pubKeyMatch {
		log.Fatalf("Public key in timestamper certificate does not match public key in provided private key")
	}

	if caChainFile, err := os.ReadFile(_ca_chain_file); err == nil {
		for {
			pemCert, restData := pem.Decode(caChainFile)

			if caCert, err := x509.ParseCertificate(pemCert.Bytes); err == nil {
				if caCert.Equal(_cert) {
					log.Fatal("Timestamper certificate shall not be in CAs chain")
				}

				_ca_chain = append(_ca_chain, caCert)
			} else {
				log.Fatalf("Cannot parse CA certificate: %s", err)
			}

			if len(_ca_chain) > 1 {
				prevCa := _ca_chain[len(_ca_chain)-2]
				currCa := _ca_chain[len(_ca_chain)-1]
				if err := prevCa.CheckSignatureFrom(currCa); err != nil {
					log.Fatalf("Invalid CA chain: %s (%s; parent: %s)", err, prevCa.Subject, currCa.Subject)
				}
			}

			if len(restData) == 0 {
				break
			}
			caChainFile = restData
		}
	} else {
		log.Fatalf("Cannot read CA chain: %s", err)
	}

	if len(_ca_chain) == 0 {
		log.Fatal("Invalid CA chain: No CAs in chain")
	}

	if err := _cert.CheckSignatureFrom(_ca_chain[0]); err != nil {
		log.Fatalf("Timestamper certificate not been issued by CA in the chain (%s)", err)
	}
}

func main() {
	service.Serve(
		fmt.Sprintf("%s:%d", _address, _port),
		_ca_chain, _cert, _priv,
	)
}
