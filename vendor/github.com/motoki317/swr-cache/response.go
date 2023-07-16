package swrcache

import (
	"bytes"
	"net/http"
)

type response struct {
	headers http.Header
	status  int
	body    bytes.Buffer
}

func newResponse() *response {
	return &response{
		headers: make(http.Header),
	}
}

func (r *response) Header() http.Header {
	return r.headers
}

func (r *response) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *response) WriteHeader(statusCode int) {
	r.status = statusCode
}
