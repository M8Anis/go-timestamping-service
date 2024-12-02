package timestamper

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"log"
	"strings"

	cms "github.com/github/smimesign/ietf-cms"
)

type AuthenticodeTimestampRequest struct {
	CounterSignatureType asn1.ObjectIdentifier

	ContentInfo struct {
		ContentType asn1.ObjectIdentifier

		Content struct {
			Bytes []byte // Signature
		} `asn1:"tag:0"`
	}
}

// The signature algorithm is determined from the certificate, I think
func (stamper *Timestamper) Authenticode(req []byte) (resp []byte, e *HttpError) {
	// Windows sends a nul-terminated string and disrupts the Base64 decoder in Golang)
	pemReq := strings.ReplaceAll(string(req), "\x00", "")

	// I do not use `pem.Decode`, because is no header and footer in the request
	derReq, err := base64.StdEncoding.DecodeString(pemReq)
	if err != nil {
		log.Printf("Request cannot be decoded (Base64): %s", err)
		return nil, ErrorWhileParsingRequest
	}

	tsReq := AuthenticodeTimestampRequest{}
	_, err = asn1.Unmarshal(derReq, &tsReq)
	if err != nil {
		log.Printf("Request cannot be decoded (ASN.1): %s", err)
		return nil, ErrorWhileParsingRequest
	}

	derResp, err := cms.Sign(tsReq.ContentInfo.Content.Bytes, stamper.FullChain, stamper.PrivateKey)
	if err != nil {
		log.Printf("Response cannot be signed: %s", err)
		return nil, GenericError
	}

	pemResp := pem.EncodeToMemory(&pem.Block{
		Type:  "CMS",
		Bytes: derResp,
	})
	if pemResp == nil {
		log.Println("Response encoded in PEM is NULL")
		return nil, GenericError
	}

	// Removing the PEM header and footer from the response, as in the request
	resp = pemResp[20 : len(pemResp)-18]

	return
}
