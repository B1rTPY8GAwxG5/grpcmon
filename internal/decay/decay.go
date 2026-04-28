// Package decay provides time-based score decay for captured gRPC entries.
// Scores decrease exponentially as entries age, making recent traffic more
// significant in analysis and prioritisation.
package decay

import (
	"math"
	"sync"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Options configures the decay function.
type Options struct {
	// HalfLife is the duration after which a score is halved.
	HalfLife time.Duration
}

// DefaultOptions returns sensible decay defaults.
func DefaultOptions() Options {
	return Options{
		HalfLife: 5 * time.Minute,
	}
}

// Scorer applies exponential decay to entry scores.
type Scorer struct {
	mu   sync.Mutex
	opts Options
	now  func() time.Time
}

// New creates a Scorer with the given options.
func New(opts Options) *Scorer {
	if opts.HalfLife <= 0 {
		opts.HalfLife = DefaultOptions().HalfLife
	}
	return &Scorer{
		opts: opts,
		now:  time.Now,
	}
}

// Score returns a decay-weighted score in [0, 1] for the entry.
// A score of 1 means the entry was captured right now; it halves every HalfLife.
func (s *Scorer) Score(e capture.Entry) float64 {
	s.mu.Lock()
	now := s.now()
	s.mu.Unlock()

	if e.Timestamp.IsZero() {
		return 0
	}
	age := now.Sub(e.Timestamp)
	if age < 0 {
		age = 0
	}
	hl := s.opts.HalfLife.Seconds()
	return math.Pow(0.5, age.Seconds()/hl)
}

// Apply returns entries sorted by their decay score descending.
// Entries with a score below threshold are dropped.
func (s *Scorer) Apply(entries []capture.Entry, threshold float64) []capture.Entry {
	type scored struct {
		e     capture.Entry
		score float64
	}

	var kept []scored
	for _, e := range entries {
		sc := s.Score(e)
		if sc >= threshold {
			kept = append(kept, scored{e, sc})
		}
	}

	// insertion sort — entry counts are typically small
	for i := 1; i < len(kept); i++ {
		for j := i; j > 0 && kept[j].score > kept[j-1].score; j-- {
			kept[j], kept[j-1] = kept[j-1], kept[j]
		}
	}

	out := make([]capture.Entry, len(kept))
	for i, k := range kept {
		out[i] = k.e
	}
	return out
}
