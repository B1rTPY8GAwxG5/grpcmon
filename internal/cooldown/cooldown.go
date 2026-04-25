// Package cooldown enforces a minimum interval between successive replay
// attempts for the same gRPC method, preventing thundering-herd problems
// during development replay sessions.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last-seen timestamp per method and reports whether
// enough time has elapsed to allow a new attempt.
type Cooldown struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
}

// New creates a Cooldown that enforces the given minimum interval between
// calls to the same method. If interval is zero or negative it defaults to
// one second.
func New(interval time.Duration) *Cooldown {
	if interval <= 0 {
		interval = time.Second
	}
	return &Cooldown{
		last:     make(map[string]time.Time),
		interval: interval,
	}
}

// Allow returns true and records the current time when the method has not
// been seen within the cooldown interval. It returns false without updating
// the timestamp when the interval has not yet elapsed.
func (c *Cooldown) Allow(method string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if t, ok := c.last[method]; ok && now.Sub(t) < c.interval {
		return false
	}
	c.last[method] = now
	return true
}

// Reset clears the recorded timestamp for the given method so that the next
// call to Allow will succeed immediately.
func (c *Cooldown) Reset(method string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.last, method)
}

// ResetAll clears all recorded timestamps.
func (c *Cooldown) ResetAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last = make(map[string]time.Time)
}

// Remaining returns the duration left before the given method is allowed
// again. It returns zero when the method is already allowed.
func (c *Cooldown) Remaining(method string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	t, ok := c.last[method]
	if !ok {
		return 0
	}
	elapsed := time.Since(t)
	if elapsed >= c.interval {
		return 0
	}
	return c.interval - elapsed
}
