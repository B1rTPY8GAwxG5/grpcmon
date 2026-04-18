// Package bookmark provides named bookmarks for captured gRPC entries.
package bookmark

import (
	"errors"
	"sync"

	"github.com/grpcmon/internal/capture"
)

// ErrNotFound is returned when a bookmark name does not exist.
var ErrNotFound = errors.New("bookmark: not found")

// ErrDuplicate is returned when a bookmark name is already in use.
var ErrDuplicate = errors.New("bookmark: name already exists")

// Store maps named bookmarks to capture entry IDs.
type Store struct {
	mu    sync.RWMutex
	names map[string]string // name -> entry ID
}

// New returns an empty bookmark Store.
func New() *Store {
	return &Store{names: make(map[string]string)}
}

// Add associates name with the given entry. Returns ErrDuplicate if name exists.
func (s *Store) Add(name string, e capture.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.names[name]; ok {
		return ErrDuplicate
	}
	s.names[name] = e.ID
	return nil
}

// Remove deletes a bookmark by name. Returns ErrNotFound if absent.
func (s *Store) Remove(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.names[name]; !ok {
		return ErrNotFound
	}
	delete(s.names, name)
	return nil
}

// Get returns the entry ID associated with name.
func (s *Store) Get(name string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.names[name]
	if !ok {
		return "", ErrNotFound
	}
	return id, nil
}

// List returns all bookmark names in undefined order.
func (s *Store) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.names))
	for n := range s.names {
		out = append(out, n)
	}
	return out
}
