package goretryhttp

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Request wraps the metadata to create HTTP requests.
type Request struct {
	body io.ReadSeeker
	*http.Request
}

// NewRequest creates a new wrapped HTTP request.
func NewRequest(method, url string, body io.ReadSeeker) (*Request, error) {
	var rcBody io.ReadCloser
	// wrap the body in a noop ReadCloser if non-nil. This prevents the reader from being closed by the http client.
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}

	// Make the request with the noop-closer for the body.
	httpReq, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}

	return &Request{body, httpReq}, nil
}
