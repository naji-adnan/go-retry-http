package goretryhttp

import (
	"math"
	"time"
)

// Backoff specifies a policy for how long to wait between replies
type Backoff func(min, max time.Duration, attemptNum int) time.Duration

// DefaultBackoff will perform exponential backoff based on the attempt number and limited to the min and max provided
func DefaultBackoff(min, max time.Duration, attemptNum int) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}

	return sleep
}
