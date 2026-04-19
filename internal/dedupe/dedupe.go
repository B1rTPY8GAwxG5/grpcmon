// Package dedupe provides deduplication of captured gRPC entries
// based on method and request payload similarity.
package dedupe

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/grpcmon/internal/capture"
)

// Store tracks seen entry fingerprints and filters duplicates.
type Store struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// New returns an initialised Store.
func New() *Store {
	return &Store{seen: make(map[string]struct{})}
}

// fingerprint produces a stable hash from an entry's method and request.
func fingerprint(e capture.Entry) string {
	h := sha256.New()
	h.Write([]byte(e.Method))
	h.Write([]byte(e.Request))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// IsDuplicate reports whether the entry has been seen before.
// If not seen, it records the fingerprint and returns false.
func (s *Store) IsDuplicate(e capture.Entry) bool {
	fp := fingerprint(e)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.seen[fp]; ok {
		return true
	}
	s.seen[fp] = struct{}{}
	return false
}

// Filter returns only entries that have not been seen before,
// recording each as it is encountered.
func (s *Store) Filter(entries []capture.Entry) []capture.Entry {
	out := make([]capture.Entry, 0, len(entries))
	for _, e := range entries {
		if !s.IsDuplicate(e) {
			out = append(out, e)
		}
	}
	return out
}

// Reset clears all recorded fingerprints.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = make(map[string]struct{})
}

// Len returns the number of unique fingerprints stored.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.seen)
}
