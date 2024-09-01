package main

import (
	"log"
	"net/http"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/service"
)

func main() {
	http.HandleFunc("/", service.HttpEndpoint)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
