package goretryhttp

import (
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
