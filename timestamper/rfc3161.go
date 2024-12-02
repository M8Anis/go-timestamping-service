package timestamper

import (
	"crypto"
	"encoding/asn1"
	"net/http"
	"time"

	"github.com/digitorus/timestamp"
	"github.com/sirupsen/logrus"
)

// At now, only SHA256 response signing is supported
func (stamper *Timestamper) Rfc3161(req []byte) (resp []byte, status int) {
	tsReq, err := timestamp.ParseRequest(req)
	if err != nil {
		logrus.Infof("Request cannot be parsed: %s", err)
		return nil, http.StatusBadRequest
	}

	tsResp := timestamp.Timestamp{
		Certificates:      stamper.CaChain,
		AddTSACertificate: true,

		HashAlgorithm: tsReq.HashAlgorithm,
		HashedMessage: tsReq.HashedMessage,
		Nonce:         tsReq.Nonce,

		Time: time.Now().UTC(),

		// idk but its needed
		Policy: asn1.ObjectIdentifier{0, 0, 0},
	}

	resp, err = tsResp.CreateResponseWithOpts(stamper.Certificate, stamper.PrivateKey, crypto.SHA256)
	if err != nil {
		logrus.Errorf("Response cannot be created: %s", err)
		return nil, http.StatusInternalServerError
	}

	return
}
