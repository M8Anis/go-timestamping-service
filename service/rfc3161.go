package service

import (
	"crypto"
	"encoding/asn1"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/digitorus/timestamp"
)

// At now, only SHA256 response signing is supported
func Rfc3161(w http.ResponseWriter, derReq []byte) {
	req, err := timestamp.ParseRequest(derReq)
	if err != nil {
		ErrorPage(w, http.StatusBadRequest,
			"Error while parsing the request. For more information, refer to the console",
		)
		log.Printf("Request cannot be parsed: %s", err)
		return
	}

	ts := timestamp.Timestamp{
		HashAlgorithm: req.HashAlgorithm,
		HashedMessage: req.HashedMessage,

		Nonce: req.Nonce,

		Time:     time.Now().UTC(),
		Accuracy: time.Millisecond * 100,

		AddTSACertificate: true,

		Qualified: true,
		Ordering:  true,
		Policy:    asn1.ObjectIdentifier{2, 4, 5, 6},
	}

	derResp, err := ts.CreateResponseWithOpts(signingCertificate, signingKey, crypto.SHA256)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Printf("Response cannot be created: %s", err)
		return
	}

	w.Header().Add("Content-Type", RFC3161_REPLY)
	fmt.Fprintf(w, "%s", derResp)
}
