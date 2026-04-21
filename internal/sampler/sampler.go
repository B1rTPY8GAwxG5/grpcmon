// Package sampler provides probabilistic sampling of captured gRPC entries.
// It allows reducing the volume of stored traffic by retaining only a
// configurable fraction of entries, chosen uniformly at random.
package sampler

import (
	"math/rand"
	"sync"

	"github.com/grpcmon/internal/capture"
)

// Sampler decides whether an entry should be retained based on a sampling rate.
type Sampler struct {
	mu   sync.Mutex
	rate float64 // 0.0 – 1.0
	rng  *rand.Rand
}

// New creates a Sampler with the given rate.
// rate must be in the range [0.0, 1.0]; values outside this range are clamped.
func New(rate float64, src rand.Source) *Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	if src == nil {
		src = rand.NewSource(42)
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(src), //nolint:gosec
	}
}

// Keep returns true if the entry should be retained according to the
// configured sampling rate.
func (s *Sampler) Keep(_ capture.Entry) bool {
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}

// Filter returns the subset of entries that pass the sampler.
func (s *Sampler) Filter(entries []capture.Entry) []capture.Entry {
	out := make([]capture.Entry, 0, len(entries))
	for _, e := range entries {
		if s.Keep(e) {
			out = append(out, e)
		}
	}
	return out
}

// Rate returns the current sampling rate.
func (s *Sampler) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}

// SetRate updates the sampling rate at runtime.
// Values are clamped to [0.0, 1.0].
func (s *Sampler) SetRate(rate float64) {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	s.mu.Lock()
	s.rate = rate
	s.mu.Unlock()
}
