// Package truncate provides utilities for capping a capture store to a
// maximum number of entries, evicting the oldest records first.
package truncate

import (
	"sync"

	"github.com/grpcmon/internal/capture"
)

// Truncator trims a Store to a configured maximum size.
type Truncator struct {
	mu      sync.Mutex
	store   *capture.Store
	maxSize int
}

// New returns a Truncator that will keep at most maxSize entries in store.
// If maxSize is less than 1 it defaults to 1.
func New(store *capture.Store, maxSize int) *Truncator {
	if maxSize < 1 {
		maxSize = 1
	}
	return &Truncator{store: store, maxSize: maxSize}
}

// Trim removes the oldest entries from the store until its length is at or
// below the configured maximum. It returns the number of entries removed.
func (t *Truncator) Trim() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	entries := t.store.List()
	excess := len(entries) - t.maxSize
	if excess <= 0 {
		return 0
	}

	// Entries are stored oldest-first; clear all then re-add the tail.
	keep := entries[excess:]
	t.store.Clear()
	for _, e := range keep {
		t.store.Add(e)
	}
	return excess
}

// MaxSize returns the configured maximum number of entries.
func (t *Truncator) MaxSize() int {
	return t.maxSize
}
