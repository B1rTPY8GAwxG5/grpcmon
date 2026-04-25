// Package coalesce merges multiple capture entries that share the same
// method and status code into a single representative entry, keeping the
// most-recent timestamp and averaging latency across the group.
package coalesce

import (
	"sync"
	"time"

	"github.com/example/grpcmon/internal/capture"
)

// key identifies a group of entries that can be merged.
type key struct {
	Method string
	Status uint32
}

// Coalescer merges entries that share the same method and status code.
type Coalescer struct {
	mu      sync.Mutex
	groups  map[key][]capture.Entry
}

// New returns an initialised Coalescer.
func New() *Coalescer {
	return &Coalescer{
		groups: make(map[key][]capture.Entry),
	}
}

// Add records an entry for later merging.
func (c *Coalescer) Add(e capture.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k := key{Method: e.Method, Status: e.StatusCode}
	c.groups[k] = append(c.groups[k], e)
}

// Flush merges all accumulated entries and returns one representative
// entry per (method, status) group. The internal state is reset.
func (c *Coalescer) Flush() []capture.Entry {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make([]capture.Entry, 0, len(c.groups))
	for _, entries := range c.groups {
		result = append(result, merge(entries))
	}
	c.groups = make(map[key][]capture.Entry)
	return result
}

// merge combines a slice of entries into one, averaging latency and
// keeping the latest timestamp.
func merge(entries []capture.Entry) capture.Entry {
	if len(entries) == 0 {
		return capture.Entry{}
	}
	base := entries[0]
	var totalLatency time.Duration
	for _, e := range entries {
		totalLatency += e.Latency
		if e.Timestamp.After(base.Timestamp) {
			base.Timestamp = e.Timestamp
		}
	}
	base.Latency = totalLatency / time.Duration(len(entries))
	return base
}
