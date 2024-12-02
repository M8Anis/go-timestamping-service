package service

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitea.m8anis.internal/M8Anis/go-timestamping-service/timestamper"
	"github.com/sirupsen/logrus"
)

func handleQuery(w http.ResponseWriter, r *http.Request) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if !timestamper.QueryValid(contentType) {
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
		fmt.Fprint(w, "Query must be present")
		return
	}

	resp, status := instance.MakeReply(contentType, req)
	if status == http.StatusOK {
		timestamper.AddReplyContentType(w, contentType)
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", resp)
}
