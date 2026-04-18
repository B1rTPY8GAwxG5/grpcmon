// Package middleware provides gRPC client interceptors that transparently
// capture unary calls into a capture.Store for inspection, replay and export.
//
// Usage:
//
//	store := capture.NewStore(500)
//	conn, err := middleware.DialWithCapture("localhost:50051", store,
//	    grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//
// All unary RPCs made through the returned connection are recorded with their
// method name, request/response payloads, status code and latency.
package middleware
