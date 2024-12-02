package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/timestamper"
)

const RFC3161_REPLY string = "application/timestamp-reply"
const RFC3161_QUERY string = "application/timestamp-query"

const AUTHENTICODE string = "application/octet-stream"

func HttpEndpoint(w http.ResponseWriter, r *http.Request) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if len(contentType) == 0 {
		HomePage(w)
		return
	}

	if RFC3161_QUERY != contentType && AUTHENTICODE != contentType {
		ErrorPage(w, http.StatusBadRequest,
			fmt.Sprintf(
				"`Content-Type` must be `%s` for RFC3161 or `%s` for Authenticode(tm)",
				RFC3161_QUERY, AUTHENTICODE,
			),
		)
		return
	}

	if http.MethodPost != r.Method {
		ErrorPage(w, http.StatusMethodNotAllowed,
			fmt.Sprintf(
				"Method `%s` is not allowed", r.Method,
			),
		)
		return
	}

	defer r.Body.Close()
	req, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Body cannot be read: %s", err)
		ErrorPage(w, timestamper.GenericError.Code(), timestamper.GenericError.Error())
		return
	}
	if len(req) == 0 {
		ErrorPage(w, http.StatusBadRequest, "Request body must be present")
		return
	}

	var e *timestamper.HttpError
	var resp []byte
	switch contentType {
	case RFC3161_QUERY:
		if resp, e = instance.Rfc3161(req); e != nil {
			ErrorPage(w, e.Code(), e.Error())
			return
		}
		w.Header().Add("Content-Type", RFC3161_REPLY)
	case AUTHENTICODE:
		if resp, e = instance.Authenticode(req); e != nil {
			ErrorPage(w, e.Code(), e.Error())
			return
		}
		w.Header().Add("Content-Type", AUTHENTICODE)
	}

	fmt.Fprintf(w, "%s", resp)
}
