// Package group provides grouping of captured entries by a key function.
package group

import (
	"sync"

	"github.com/grpcmon/internal/capture"
)

// KeyFunc extracts a grouping key from an entry.
type KeyFunc func(e capture.Entry) string

// ByMethod groups entries by their gRPC method name.
func ByMethod(e capture.Entry) string { return e.Method }

// ByStatus groups entries by their status code string.
func ByStatus(e capture.Entry) string { return e.Status.String() }

// Group holds entries sharing the same key.
type Group struct {
	Key     string
	Entries []capture.Entry
}

// Store groups entries from a capture store using a key function.
type Store struct {
	mu  sync.RWMutex
	key KeyFunc
}

// New returns a Store that groups entries using fn.
func New(fn KeyFunc) *Store {
	if fn == nil {
		fn = ByMethod
	}
	return &Store{key: fn}
}

// Apply groups the provided entries and returns a slice of Group values
// sorted by insertion order of first occurrence.
func (s *Store) Apply(entries []capture.Entry) []Group {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order := make([]string, 0)
	index := make(map[string]int)
	groups := make([]Group, 0)

	for _, e := range entries {
		k := s.key(e)
		if i, ok := index[k]; ok {
			groups[i].Entries = append(groups[i].Entries, e)
		} else {
			index[k] = len(groups)
			order = append(order, k)
			groups = append(groups, Group{Key: k, Entries: []capture.Entry{e}})
		}
	}
	_ = order
	return groups
}
