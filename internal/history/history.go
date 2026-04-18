// Package history tracks replay history, recording which entries have been
// replayed and their outcomes for later inspection.
package history

import (
	"sync"
	"time"

	"grpcmon/internal/capture"
)

// Record holds the result of a single replay attempt.
type Record struct {
	EntryID   string
	Method    string
	ReplayedAt time.Time
	Success   bool
	Error     string
}

// History stores replay records in insertion order with a configurable cap.
type History struct {
	mu      sync.RWMutex
	records []Record
	maxSize int
}

// New returns a History that retains at most maxSize records.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{maxSize: maxSize}
}

// Add appends a replay record, evicting the oldest when full.
func (h *History) Add(e capture.Entry, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := Record{
		EntryID:    e.ID,
		Method:     e.Method,
		ReplayedAt: time.Now().UTC(),
		Success:    err == nil,
	}
	if err != nil {
		r.Error = err.Error()
	}
	if len(h.records) >= h.maxSize {
		h.records = h.records[1:]
	}
	h.records = append(h.records, r)
}

// List returns a shallow copy of all records.
func (h *History) List() []Record {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Record, len(h.records))
	copy(out, h.records)
	return out
}

// Clear removes all records.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = h.records[:0]
}
