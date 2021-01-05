package goretryhttp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

var (
	// Default retry configuration
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
	// defaultClient       = NewClient()
	respReadLimit = int64(4096)
)

// Client is a wrapper of the standard http.Client but has additional functionality for automatic retries
type Client struct {
	HTTPClient    *http.Client  // Standard HTTP Client
	RetryWaitMin  time.Duration // Min time to wait
	RetryWaitMax  time.Duration // Max time to wait
	RetryMax      int           // Max number of retries
	CheckForRetry CheckForRetry // specifies the policy for handling retries and is called after each request
	Backoff       Backoff       // specifies the policy for how long to wait between retries
}

// NewClient created a new Client
func NewClient() *Client {
	return &Client{
		HTTPClient:    cleanhttp.DefaultClient(),
		RetryWaitMin:  defaultRetryWaitMin,
		RetryWaitMax:  defaultRetryWaitMax,
		RetryMax:      defaultRetryMax,
		CheckForRetry: DefaultRetryPolicy,
		Backoff:       DefaultBackoff,
	}
}

// Do wraps calling an HTTP method with retries.
func (c *Client) Do(req *Request) (*http.Response, error) {
	log.Printf("[DEBUG] %s %s", req.Method, req.URL)
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
		shallRetry, checkErr := c.CheckForRetry(resp, err)
		if err != nil {
			log.Printf("[Err] %s %s request failed: %v", req.Method, req.URL, err)
		}

		// Decide if we should continue
		if !shallRetry {
			if checkErr != nil {
				err = checkErr
			}
			return resp, err
		}

		// We're going to retry, consume any response to reuse the connection.
		// if err == nil {
		// 	c.drainBodyOnFailResponse(resp.Body)
		// }

		remain := c.RetryMax - i
		if remain == 0 {
			break
		}
		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}
		log.Printf("[DEBUG] %s: retrying in %s (%d left)", desc, wait, remain)
		time.Sleep(wait)
	}

	// Return an error if we fall out of the retry loop
	return nil, fmt.Errorf("%s %s giving up after %d attempts", req.Method, req.URL, c.RetryMax+1)

}

// // Try to read the response body so we can reuse this connection
// func (c *Client) drainBodyOnFailResponse(body io.ReadCloser) {
// 	defer body.Close()

// 	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
// 	if err != nil {
// 		fmt.Printf("[ERR] error reading response body: %v", err)
// 	}
// }

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
