package service

import (
	"net/http"
)

func Rfc3161(w http.ResponseWriter, body []byte) {
	ErrorPage(w, http.StatusInternalServerError, "RFC3161 not implemented")
}
