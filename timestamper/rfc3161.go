package timestamper

import (
	"crypto"
	"encoding/asn1"
	"log"
	"time"

	"github.com/digitorus/timestamp"
)

// At now, only SHA256 response signing is supported
func (stamper *Timestamper) Rfc3161(req []byte) (resp []byte, e *HttpError) {
	tsReq, err := timestamp.ParseRequest(req)
	if err != nil {
		log.Printf("Request cannot be parsed: %s", err)
		return nil, ErrorWhileParsingRequest
	}

	tsResp := timestamp.Timestamp{
		AddTSACertificate: true,

		HashAlgorithm: tsReq.HashAlgorithm,
		HashedMessage: tsReq.HashedMessage,
		Nonce:         tsReq.Nonce,

		Time: time.Now().UTC(),

		// idk but its needed
		Policy: asn1.ObjectIdentifier{0, 0, 0},
	}
	if stamper.ChainLength > 1 {
		tsResp.Certificates = stamper.CertChain
	}

	resp, err = tsResp.CreateResponseWithOpts(stamper.Certificate, stamper.PrivateKey, crypto.SHA256)
	if err != nil {
		log.Printf("Response cannot be created: %s", err)
		return nil, GenericError
	}

	return
}
