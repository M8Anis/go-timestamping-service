package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const RFC3161_CONTENT_TYPE string = "application/timestamp-query"
const AUTHENTICODE_CONTENT_TYPE string = "application/octet-stream"

func HttpEndpoint(w http.ResponseWriter, r *http.Request) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if len(contentType) == 0 {
		HomePage(w)
		return
	}

	if RFC3161_CONTENT_TYPE != contentType && AUTHENTICODE_CONTENT_TYPE != contentType {
		ErrorPage(w, http.StatusBadRequest,
			fmt.Sprintf(
				"`Content-Type` must be `%s` for RFC3161 or `%s` for Authenticode(tm)",
				RFC3161_CONTENT_TYPE, AUTHENTICODE_CONTENT_TYPE,
			),
		)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError,
			"Error has occured. For more information, refer to the console",
		)
		log.Printf("Body cannot be read: %s", err)
		return
	}
	if len(body) == 0 {
		ErrorPage(w, http.StatusBadRequest, "Request body must be present")
		return
	}

	switch contentType {
	case "application/timestamp-query":
		Rfc3161(w, body)
	case "application/octet-stream":
		Authenticode(w, body)
	}
}
