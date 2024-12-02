package service

import (
	"crypto"
	"crypto/x509"
	"log"
	"net/http"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/timestamper"
)

var instance *timestamper.Timestamper

func Serve(certChainLen int, fullCertChain, certChain []*x509.Certificate, timestamperCert *x509.Certificate, timestamperPrivKey crypto.Signer) {
	instance = &timestamper.Timestamper{
		CertChain:   certChain,
		Certificate: timestamperCert,
		PrivateKey:  timestamperPrivKey,

		FullCertChain: fullCertChain,
		ChainLength:   certChainLen,
	}

	http.HandleFunc("/", HttpEndpoint)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
