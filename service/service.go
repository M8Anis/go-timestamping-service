package service

import (
	"context"
	"crypto"
	"crypto/x509"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/timestamper"
	"github.com/gorilla/mux"
)

var instance *timestamper.Timestamper

func Serve(host string, caChain []*x509.Certificate, stamperCert *x509.Certificate, stamperPrivKey crypto.Signer) {
	r := mux.NewRouter()
	r.HandleFunc("/", handleQuery).
		Methods(http.MethodPost)

	instance = &timestamper.Timestamper{
		CaChain:     caChain,
		Certificate: stamperCert,
		PrivateKey:  stamperPrivKey,

		FullChain: append([]*x509.Certificate{stamperCert}, caChain...),

		Server: &http.Server{
			Handler: r,
			Addr:    host,
		},
	}

	go func() {
		instance.Server.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	instance.Server.Shutdown(ctx)
}
