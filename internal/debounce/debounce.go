// Package debounce provides a mechanism for coalescing rapid bursts of
// capture entries into a single notification, reducing downstream noise
// during high-frequency gRPC traffic.
package debounce

import (
	"context"
	"sync"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Handler is called with the accumulated entries after the debounce window
// has elapsed without a new entry arriving.
type Handler func(entries []capture.Entry)

// Debouncer coalesces rapid entry additions into batched notifications.
type Debouncer struct {
	mu       sync.Mutex
	window   time.Duration
	handler  Handler
	pending  []capture.Entry
	timer    *time.Timer
}

// New creates a Debouncer that waits for window of silence before invoking
// handler with all accumulated entries. A zero or negative window defaults
// to 100ms.
func New(window time.Duration, handler Handler) *Debouncer {
	if window <= 0 {
		window = 100 * time.Millisecond
	}
	return &Debouncer{
		window:  window,
		handler: handler,
	}
}

// Add queues an entry and resets the debounce timer.
func (d *Debouncer) Add(e capture.Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.pending = append(d.pending, e)

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.window, d.flush)
}

// Flush immediately fires the handler with any pending entries, regardless
// of whether the debounce window has elapsed. Safe to call concurrently.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
	d.mu.Unlock()
	d.flush()
}

// Run starts a background goroutine that calls Flush when ctx is cancelled,
// ensuring no entries are lost on shutdown.
func (d *Debouncer) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		d.Flush()
	}()
}

func (d *Debouncer) flush() {
	d.mu.Lock()
	if len(d.pending) == 0 {
		d.mu.Unlock()
		return
	}
	batch := make([]capture.Entry, len(d.pending))
	copy(batch, d.pending)
	d.pending = d.pending[:0]
	d.mu.Unlock()

	d.handler(batch)
}
