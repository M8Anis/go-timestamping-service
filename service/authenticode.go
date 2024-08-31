package service

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	cms "github.com/github/smimesign/ietf-cms"
)

type AuthenticodeTimestampRequest struct {
	OID asn1.ObjectIdentifier

	Payload struct {
		OID asn1.ObjectIdentifier

		Data asn1.RawValue
	}
}

func Authenticode(w http.ResponseWriter, body []byte) {
	rawAsnData := make([]byte, 128)
	n, err := base64.StdEncoding.Decode(rawAsnData, body)
	if err != nil && n != len(body)-1 {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Printf("Body cannot be decoded (Base64): %s", err)
		return
	}
	rawAsnData = rawAsnData[:n]

	req := &AuthenticodeTimestampRequest{}
	_, err = asn1.Unmarshal(rawAsnData, req)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Printf("Body cannot be decoded (ASN.1): %s", err)
		return
	}

	signedData, err := cms.Sign(req.Payload.Data.Bytes[2:], []*x509.Certificate{signingCertificate}, signingKey)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Printf("Body cannot be signed: %s", err)
		return
	}

	pemSigData := pem.EncodeToMemory(&pem.Block{
		Type:  "CMS",
		Bytes: signedData,
	})
	if pemSigData == nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Print("Signed data PEM encoded equals NULL")
		return
	}
	pemSigData = pemSigData[20 : len(pemSigData)-18]

	w.Header().Add("Content-Type", AUTHENTICODE_CONTENT_TYPE)
	fmt.Fprintf(w, "%s", pemSigData)
}
