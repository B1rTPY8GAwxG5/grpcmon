// Package timeout provides per-method deadline enforcement for replayed gRPC entries.
// Each method can be assigned an independent timeout; a default fallback is used
// when no method-specific value is registered.
package timeout

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/grpcmon/internal/capture"
)

// ErrDeadlineExceeded is returned when a replayer does not complete within the
// allowed duration.
var ErrDeadlineExceeded = errors.New("timeout: deadline exceeded")

// Replayer is the function signature used to replay a single capture entry.
type Replayer func(ctx context.Context, e capture.Entry) error

// Manager holds per-method and default timeout values.
type Manager struct {
	mu         sync.RWMutex
	defaultTTL time.Duration
	methods    map[string]time.Duration
}

// New returns a Manager with the given default timeout.
// If d is zero or negative it is clamped to 5 seconds.
func New(d time.Duration) *Manager {
	if d <= 0 {
		d = 5 * time.Second
	}
	return &Manager{
		defaultTTL: d,
		methods:    make(map[string]time.Duration),
	}
}

// Set registers a timeout for a specific gRPC method.
func (m *Manager) Set(method string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.methods[method] = d
}

// Get returns the timeout for the given method, falling back to the default.
func (m *Manager) Get(method string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if d, ok := m.methods[method]; ok {
		return d
	}
	return m.defaultTTL
}

// Wrap returns a Replayer that enforces the method-specific (or default)
// deadline before delegating to next.
func (m *Manager) Wrap(next Replayer) Replayer {
	return func(ctx context.Context, e capture.Entry) error {
		timeout := m.Get(e.Method)
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		type result struct{ err error }
		ch := make(chan result, 1)
		go func() {
			ch <- result{err: next(ctx, e)}
		}()

		select {
		case r := <-ch:
			return r.err
		case <-ctx.Done():
			return ErrDeadlineExceeded
		}
	}
}
