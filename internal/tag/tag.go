// Package tag provides tagging and lookup for captured gRPC entries.
package tag

import (
	"sync"

	"github.com/grpcmon/internal/capture"
)

// Store maps string tags to sets of entry IDs.
type Store struct {
	mu   sync.RWMutex
	tags map[string]map[string]struct{}
}

// New returns an empty tag Store.
func New() *Store {
	return &Store{tags: make(map[string]map[string]struct{})}
}

// Add associates the given tags with the entry ID.
func (s *Store) Add(id string, tags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tags {
		if s.tags[t] == nil {
			s.tags[t] = make(map[string]struct{})
		}
		s.tags[t][id] = struct{}{}
	}
}

// Remove disassociates the given tags from the entry ID.
func (s *Store) Remove(id string, tags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tags {
		delete(s.tags[t], id)
	}
}

// Lookup returns all entry IDs associated with a tag.
func (s *Store) Lookup(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, 0, len(s.tags[tag]))
	for id := range s.tags[tag] {
		ids = append(ids, id)
	}
	return ids
}

// Filter returns entries whose IDs are tagged with the given tag.
func (s *Store) Filter(tag string, entries []capture.Entry) []capture.Entry {
	ids := s.Lookup(tag)
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	out := make([]capture.Entry, 0)
	for _, e := range entries {
		if _, ok := set[e.ID]; ok {
			out = append(out, e)
		}
	}
	return out
}

// Tags returns all known tags.
func (s *Store) Tags() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.tags))
	for t := range s.tags {
		out = append(out, t)
	}
	return out
}
