package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const RFC3161_REPLY string = "application/timestamp-query"
const RFC3161_QUERY string = "application/timestamp-reply"

const AUTHENTICODE_CONTENT_TYPE string = "application/octet-stream"

func HttpEndpoint(w http.ResponseWriter, r *http.Request) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if len(contentType) == 0 {
		HomePage(w)
		return
	}

	if RFC3161_QUERY != contentType && AUTHENTICODE_CONTENT_TYPE != contentType {
		ErrorPage(w, http.StatusBadRequest,
			fmt.Sprintf(
				"`Content-Type` must be `%s` for RFC3161 or `%s` for Authenticode(tm)",
				RFC3161_QUERY, AUTHENTICODE_CONTENT_TYPE,
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
	case RFC3161_QUERY:
		Rfc3161(w, body)
	case AUTHENTICODE_CONTENT_TYPE:
		Authenticode(w, string(body))
	}
}
