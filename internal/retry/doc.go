// Package retry provides a simple gRPC-aware retry policy for use when
// replaying or re-issuing captured requests.
//
// Use DefaultPolicy to obtain a sensible starting configuration, then call Do
// with the function that performs the gRPC call. Do will retry up to
// MaxAttempts times with a fixed Backoff delay when the returned error carries
// one of the configured retryable status codes.
package retry
