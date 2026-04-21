// Package cursor provides a simple stateful cursor for navigating
// a list of captured entries by index, with bounds-safe movement.
package cursor

import (
	"errors"
	"sync"

	"github.com/grpcmon/internal/capture"
)

// ErrEmpty is returned when the cursor operates on an empty entry list.
var ErrEmpty = errors.New("cursor: no entries")

// Cursor holds the current position within a slice of capture entries.
type Cursor struct {
	mu      sync.Mutex
	entries []capture.Entry
	pos     int
}

// New creates a Cursor initialised at the first entry.
func New(entries []capture.Entry) *Cursor {
	return &Cursor{entries: entries}
}

// Current returns the entry at the current position.
func (c *Cursor) Current() (capture.Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) == 0 {
		return capture.Entry{}, ErrEmpty
	}
	return c.entries[c.pos], nil
}

// Next advances the cursor by one and returns the new current entry.
// It does not wrap around; calling Next at the last entry is a no-op.
func (c *Cursor) Next() (capture.Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) == 0 {
		return capture.Entry{}, ErrEmpty
	}
	if c.pos < len(c.entries)-1 {
		c.pos++
	}
	return c.entries[c.pos], nil
}

// Prev moves the cursor back by one and returns the new current entry.
// It does not wrap around; calling Prev at position 0 is a no-op.
func (c *Cursor) Prev() (capture.Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) == 0 {
		return capture.Entry{}, ErrEmpty
	}
	if c.pos > 0 {
		c.pos--
	}
	return c.entries[c.pos], nil
}

// Pos returns the current zero-based index.
func (c *Cursor) Pos() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pos
}

// Len returns the total number of entries held by the cursor.
func (c *Cursor) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// Reset moves the cursor back to position 0.
func (c *Cursor) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pos = 0
}
