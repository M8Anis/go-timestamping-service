package timestamper

import (
	"crypto"
	"crypto/x509"
	"net/http"

	"github.com/sirupsen/logrus"
)

const RFC3161_REPLY string = "application/timestamp-reply"
const RFC3161_QUERY string = "application/timestamp-query"

const AUTHENTICODE string = "application/octet-stream"

type Timestamper struct {
	FullChain []*x509.Certificate

	CaChain     []*x509.Certificate
	Certificate *x509.Certificate
	PrivateKey  crypto.Signer

	Server *http.Server
}

func (stamper *Timestamper) MakeReply(contentType string, query []byte) (reply []byte, status int) {
	switch contentType {
	case RFC3161_QUERY:
		reply, status = stamper.rfc3161(query)
	case AUTHENTICODE:
		reply, status = stamper.authenticode(query)
	default:
		logrus.Fatalf("Unknown query type: %s", contentType)
	}
	return
}

func QueryValid(contentType string) bool {
	switch contentType {
	case RFC3161_QUERY:
	case AUTHENTICODE:
		return true
	default:
	}
	return false
}

func AddReplyContentType(w http.ResponseWriter, queryContentType string) {
	switch queryContentType {
	case RFC3161_QUERY:
		w.Header().Add("Content-Type", RFC3161_REPLY)
	case AUTHENTICODE:
		w.Header().Add("Content-Type", AUTHENTICODE)
	default:
		logrus.Fatalf("Unknown query type: %s", queryContentType)
	}
	return
}
