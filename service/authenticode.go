package service

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strings"

	cms "github.com/github/smimesign/ietf-cms"
)

type AuthenticodeTimestampRequest struct {
	OID asn1.ObjectIdentifier

	Payload struct {
		OID asn1.ObjectIdentifier

		// idk how to correctly parse this structure
		Data asn1.RawValue
	}
}

// The signature algorithm is determined from the certificate, I think
func Authenticode(w http.ResponseWriter, pemReq string) {
	// Windows sends a nul-terminated string and disrupts the Base64 decoder in Golang)
	pemReq = strings.ReplaceAll(pemReq, "\x00", "")

	// I do not use `pem.Decode`, because is no header and footer in the request
	derReq, err := base64.StdEncoding.DecodeString(pemReq)
	if err != nil {
		ErrorPage(w, http.StatusBadRequest,
			"Error while parsing the request. For more information, refer to the console",
		)
		log.Printf("Request cannot be decoded (Base64): %s", err)
		return
	}

	req := &AuthenticodeTimestampRequest{}
	_, err = asn1.Unmarshal(derReq, req)
	if err != nil {
		ErrorPage(w, http.StatusBadRequest,
			"Error while parsing the request. For more information, refer to the console",
		)
		log.Printf("Request cannot be decoded (ASN.1): %s", err)
		return
	}

	// Cropping the `req.Payload.Data.Bytes`, because I dunno how to parse this structure
	derResp, err := cms.Sign(req.Payload.Data.Bytes[2:], []*x509.Certificate{signingCertificate}, signingKey)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Printf("Response cannot be signed: %s", err)
		return
	}

	pemResp := pem.EncodeToMemory(&pem.Block{
		Type:  "CMS",
		Bytes: derResp,
	})
	if pemResp == nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Print("Response encoded in PEM is NULL")
		return
	}
	// Removing the PEM header and footer from the response, as in the request
	pemResp = pemResp[20 : len(pemResp)-18]

	w.Header().Add("Content-Type", AUTHENTICODE)
	fmt.Fprintf(w, "%s", pemResp)
}
