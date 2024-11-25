package http

import (
	"fmt"
	"io"
	"net/http"
)

const (
	errorResponseMaxSize = 1024
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Error struct {
	Response *http.Response
	Body     []byte
}

func NewError(res *http.Response) *Error {
	body, _ := io.ReadAll(io.LimitReader(res.Body, errorResponseMaxSize))
	return &Error{
		Response: res,
		Body:     body,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("HTTP error: %d %s\nBody: %s", e.Response.StatusCode, e.Response.Status, string(e.Body))
}
