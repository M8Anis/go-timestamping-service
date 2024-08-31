package service

import (
	"net/http"
)

func Authenticode(w http.ResponseWriter, body []byte) {
	ErrorPage(w, http.StatusInternalServerError, "Authenticode not implemented")
}
