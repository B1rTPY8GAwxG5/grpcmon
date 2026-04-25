// Package circuitbreaker implements a simple circuit breaker for gRPC replay
// operations. It tracks consecutive failures and opens the circuit after a
// configurable threshold, preventing further calls until a cooldown elapses.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuitbreaker: circuit is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Options configures the circuit breaker behaviour.
type Options struct {
	// MaxFailures is the number of consecutive failures before opening.
	MaxFailures int
	// Cooldown is the duration to wait before moving to half-open.
	Cooldown time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxFailures: 5,
		Cooldown:    10 * time.Second,
	}
}

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	opts     Options
	mu       sync.Mutex
	state    State
	failures int
	openedAt time.Time
}

// New creates a Breaker with the given options.
func New(opts Options) *Breaker {
	if opts.MaxFailures <= 0 {
		opts.MaxFailures = 1
	}
	if opts.Cooldown <= 0 {
		opts.Cooldown = time.Second
	}
	return &Breaker{opts: opts}
}

// Allow reports whether a call should be allowed through.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateOpen:
		if time.Since(b.openedAt) >= b.opts.Cooldown {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess records a successful call, closing the circuit if half-open.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed call, opening the circuit when the threshold
// is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.opts.MaxFailures {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current circuit state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
