package service

import (
	"crypto"
	"crypto/x509"
	"net/http"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/timestamper"
	"github.com/sirupsen/logrus"
)

var instance *timestamper.Timestamper

func Serve(host string, caChain []*x509.Certificate, stamperCert *x509.Certificate, stamperPrivKey crypto.Signer) {
	instance = &timestamper.Timestamper{
		CaChain:     caChain,
		Certificate: stamperCert,
		PrivateKey:  stamperPrivKey,

		FullChain: append([]*x509.Certificate{stamperCert}, caChain...),
	}

	http.HandleFunc("/", HttpEndpoint)

	logrus.Fatal(http.ListenAndServe(host, nil))
}
