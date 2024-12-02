package service

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const RFC3161_REPLY string = "application/timestamp-reply"
const RFC3161_QUERY string = "application/timestamp-query"

const AUTHENTICODE string = "application/octet-stream"

func handleQuery(w http.ResponseWriter, r *http.Request) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if RFC3161_QUERY != contentType && AUTHENTICODE != contentType {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	req, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Body cannot be read: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(req) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body must be present")
		return
	}

	var status int
	var resp []byte
	switch contentType {
	case RFC3161_QUERY:
		if resp, status = instance.Rfc3161(req); status != 0 {
			w.WriteHeader(status)
			return
		}
		w.Header().Add("Content-Type", RFC3161_REPLY)
	case AUTHENTICODE:
		if resp, status = instance.Authenticode(req); status != 0 {
			w.WriteHeader(status)
			return
		}
		w.Header().Add("Content-Type", AUTHENTICODE)
	}

	fmt.Fprintf(w, "%s", resp)
}
