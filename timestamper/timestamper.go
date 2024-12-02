package timestamper

import (
	"crypto"
	"crypto/x509"
	"net/http"
)

type Timestamper struct {
	FullChain []*x509.Certificate

	CaChain     []*x509.Certificate
	Certificate *x509.Certificate
	PrivateKey  crypto.Signer

	Server *http.Server
}
