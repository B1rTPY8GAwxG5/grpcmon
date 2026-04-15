package capture

import "time"

// Entry represents a single captured gRPC call.
type Entry struct {
	// ID is a unique identifier for this entry.
	ID string

	// Method is the full gRPC method path, e.g. "/package.Service/Method".
	Method string

	// StatusCode is the gRPC status code string, e.g. "OK", "NOT_FOUND".
	StatusCode string

	// Duration is the elapsed time of the RPC call.
	Duration time.Duration

	// Timestamp is when the call was captured.
	Timestamp time.Time

	// RequestBody holds the raw JSON-encoded request payload.
	RequestBody []byte

	// ResponseBody holds the raw JSON-encoded response payload.
	ResponseBody []byte

	// Metadata contains the gRPC metadata headers sent with the request.
	Metadata map[string][]string
}
