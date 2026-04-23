// Package window provides a sliding time-window view over captured entries.
// It returns only entries whose timestamp falls within a configurable duration
// relative to the current time, making it easy to surface recent traffic.
package window

import (
	"sync"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Window filters entries to those captured within a rolling time range.
type Window struct {
	mu       sync.Mutex
	store    *capture.Store
	duration time.Duration
	now      func() time.Time // injectable for testing
}

// New creates a Window that retains entries within d of the current time.
// If d is zero or negative it defaults to one minute.
func New(store *capture.Store, d time.Duration) *Window {
	if d <= 0 {
		d = time.Minute
	}
	return &Window{
		store:    store,
		duration: d,
		now:      time.Now,
	}
}

// Entries returns all captured entries whose Timestamp falls within the
// window's duration ending at the current time.
func (w *Window) Entries() []capture.Entry {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := w.now().Add(-w.duration)
	all := w.store.List()
	out := make([]capture.Entry, 0, len(all))
	for _, e := range all {
		if !e.Timestamp.Before(cutoff) {
			out = append(out, e)
		}
	}
	return out
}

// Duration returns the configured window duration.
func (w *Window) Duration() time.Duration {
	return w.duration
}
