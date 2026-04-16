// Package watch provides a simple polling watcher that periodically
// snapshots the capture store and notifies subscribers of new entries.
package watch

import (
	"context"
	"time"

	"github.com/example/grpcmon/internal/capture"
)

// Handler is called whenever new entries are detected.
type Handler func(entries []capture.Entry)

// Watcher polls a capture store at a fixed interval and invokes a
// Handler when the number of entries has grown since the last check.
type Watcher struct {
	store    *capture.Store
	interval time.Duration
	handler  Handler
}

// New creates a Watcher that polls store every interval and calls h
// when new entries are available.
func New(store *capture.Store, interval time.Duration, h Handler) *Watcher {
	return &Watcher{
		store:    store,
		interval: interval,
		handler:  h,
	}
}

// Run starts the polling loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	seen := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			all := w.store.List()
			if len(all) > seen {
				w.handler(all[seen:])
				seen = len(all)
			}
		}
	}
}
