// Package annotate provides support for attaching free-form notes to captured
// gRPC entries by their ID.
package annotate

import (
	"errors"
	"sync"
)

// ErrNotFound is returned when no annotation exists for the given ID.
var ErrNotFound = errors.New("annotate: no annotation for ID")

// Store holds annotations keyed by entry ID.
type Store struct {
	mu   sync.RWMutex
	notes map[string]string
}

// New returns an initialised annotation Store.
func New() *Store {
	return &Store{notes: make(map[string]string)}
}

// Set attaches or replaces the note for the given entry ID.
func (s *Store) Set(id, note string) error {
	if id == "" {
		return errors.New("annotate: id must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[id] = note
	return nil
}

// Get retrieves the note for the given entry ID.
func (s *Store) Get(id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notes[id]
	if !ok {
		return "", ErrNotFound
	}
	return n, nil
}

// Delete removes the annotation for the given entry ID.
func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.notes, id)
}

// List returns a copy of all annotations.
func (s *Store) List() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.notes))
	for k, v := range s.notes {
		out[k] = v
	}
	return out
}
