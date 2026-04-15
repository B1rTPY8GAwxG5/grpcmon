package capture

import (
	"context"
	"sync"
	"time"
)

// Entry represents a single captured gRPC call.
type Entry struct {
	ID        string
	Timestamp time.Time
	Method    string
	Metadata  map[string][]string
	Request   []byte
	Response  []byte
	Duration  time.Duration
	Error     string
}

// Store holds captured gRPC entries in memory.
type Store struct {
	mu      sync.RWMutex
	entries []*Entry
	maxSize int
}

// NewStore creates a new Store with the given maximum capacity.
func NewStore(maxSize int) *Store {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &Store{
		entries: make([]*Entry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add appends a new entry to the store, evicting the oldest if at capacity.
func (s *Store) Add(e *Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, e)
}

// List returns a snapshot of all captured entries.
func (s *Store) List() []*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Entry, len(s.entries))
	copy(result, s.entries)
	return result
}

// Clear removes all entries from the store.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = s.entries[:0]
}

// Len returns the current number of stored entries.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// Recorder provides helpers for recording gRPC calls into a Store.
type Recorder struct {
	store *Store
}

// NewRecorder creates a Recorder backed by the given Store.
func NewRecorder(store *Store) *Recorder {
	return &Recorder{store: store}
}

// Record saves a completed gRPC call into the store.
func (r *Recorder) Record(_ context.Context, entry *Entry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	r.store.Add(entry)
}
