package service

import (
	"fmt"
	"net/http"
)

func HomePage(w http.ResponseWriter) {
	fmt.Fprintln(w, "Timestamping service!")
}

func ErrorPage(w http.ResponseWriter, status int, description string) {
	w.WriteHeader(status)
	fmt.Fprintln(w, description)
}
