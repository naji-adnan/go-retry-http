# go-retry-http
`go-retry-http` provides an http client with retry mechanism since Go's net/http client doesn’t provide a retry mechanism by default. 

This is inspired by [HTTP Retries in Go](https://medium.com/@nitishkr88/http-retries-in-go-e622e51d249f)

---

A **retry policy** is exported as function to enable applications to implement their own retry policy. This package provides a default retry policy though, which checks for a range of status codes in the response and retry on 500 range responses.

Similarly, a **backoff** policy is exported as function type and a default implementation is provided too.