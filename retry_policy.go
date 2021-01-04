package goretryhttp

import "net/http"

// CheckForRetry is a policy for handling retries.
type CheckForRetry func(resp *http.Response, err error) (bool, error)

// DefaultRetryPolicy is a default policy to retry on connection errors and server errors.
func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true, nil
	}
	return false, nil
}
