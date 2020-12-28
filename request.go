package goretryhttp

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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

// Do wraps calling an HTTP method with retries.
func (c *Client) Do(req *Request) (*http.Response, error) {
	fmt.Printf("[DEBUG] %s %s", req.Method, req.URL)
	for i := 0; ; i++ {
		var code int

		// Always rewind the request body when non-nil
		if req.body != nil {
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("failed to seek body: %v", err)
			}
		}

		// Attempt the request
		resp, err := c.HTTPClient.Do(req.Request)
		// Check if we should continue with retries.
		checkOK, checkErr := c.CheckForRetry(resp, err)
		if err != nil {
			fmt.Printf("[Err] %s %s request failed: %v", req.Method, req.URL, err)
		} else {
			fmt.Printf("[SUCCESS] %s %s request finished successfully", req.Method, req.URL)
		}

		// Decide if we should continue
		if !checkOK {
			if checkErr != nil {
				err = checkErr
			}
			return resp, err
		}

		// We're going to retry, consume any response to reuse the connection.
		if err != nil {
			c.drainBody(resp.Body)
		}

		remain := c.RetryMax - i
		if remain == 0 {
			break
		}
		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}
		fmt.Printf("[DEBUG] %s: retrying in %s (%d left)", desc, wait, remain)
		time.Sleep(wait)
	}

	// Return an error if we fall out of the retry loop
	return nil, fmt.Errorf("%s %s giving up after %d attempts", req.Method, req.URL, c.RetryMax+1)

}

// Try to read the response body so we can reuse this connection
func (c *Client) drainBody(body io.ReadCloser) {
	defer body.Close()

	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
	if err != nil {
		fmt.Printf("[ERR] error reading response body: %v", err)
	}
}

// Get is a convenience helper for doing GET requests.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post is a convenience helper for doing POST requests.
func (c *Client) Post(url, bodyType string, body io.ReadSeeker) (*http.Response, error) {
	req, err := NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return c.Do(req)
}

// Head is a convenience method for doing simple HEAD requests.
func (c *Client) Head(url string) (*http.Response, error) {
	req, err := NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
