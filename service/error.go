package service

import "net/http"

var ErrorWhileParsingRequest *HttpError = NewHttpError(
	"Error while parsing the request. For more information, refer to the console", http.StatusBadRequest,
)
var GenericError *HttpError = NewHttpError(
	"Error has occured. For more information, refer to the console", http.StatusInternalServerError,
)

func NewHttpError(text string, code int) *HttpError {
	return &HttpError{text, code}
}

type HttpError struct {
	text string
	code int
}

func (e *HttpError) Error() string {
	return e.text
}

func (e *HttpError) Code() int {
	return e.code
}
