package timestamper

import (
	"crypto"
	"crypto/x509"
)

type Timestamper struct {
	CertChain   []*x509.Certificate
	Certificate *x509.Certificate
	PrivateKey  crypto.Signer

	FullCertChain []*x509.Certificate
	ChainLength   int
}
