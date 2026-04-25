// Package circuitbreaker provides a thread-safe circuit breaker designed for
// use with gRPC replay operations in grpcmon.
//
// A Breaker starts in the closed state, allowing all calls through. After a
// configurable number of consecutive failures it transitions to the open state
// and rejects calls immediately with ErrOpen. Once the cooldown period has
// elapsed the breaker enters the half-open state: a single probe call is
// permitted. A successful probe closes the circuit; a failed probe re-opens it.
//
// Usage:
//
//	b := circuitbreaker.New(circuitbreaker.DefaultOptions())
//	protected := circuitbreaker.Wrap(b, myReplayFunc)
//
// The Wrap helper integrates the breaker with any ReplayFunc, recording
// successes and failures automatically.
package circuitbreaker
