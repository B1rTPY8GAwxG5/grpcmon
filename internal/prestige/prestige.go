// Package prestige provides entry scoring based on recency, error rate,
// and latency to surface the most "interesting" captures for review.
package prestige

import (
	"math"
	"sort"
	"time"

	"github.com/example/grpcmon/internal/capture"
)

// Score holds the computed prestige score for a single entry.
type Score struct {
	Entry   capture.Entry
	Value   float64 // higher means more interesting
}

// Options controls how each factor contributes to the final score.
type Options struct {
	// RecencyWeight scales the recency component (default 0.3).
	RecencyWeight float64
	// LatencyWeight scales the latency component (default 0.4).
	LatencyWeight float64
	// ErrorWeight scales the error component (default 0.3).
	ErrorWeight float64
	// LatencyThresholdMS is the latency (ms) above which full latency score is awarded.
	LatencyThresholdMS float64
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		RecencyWeight:      0.3,
		LatencyWeight:      0.4,
		ErrorWeight:        0.3,
		LatencyThresholdMS: 500,
	}
}

// Rank scores each entry and returns them sorted by descending prestige.
func Rank(entries []capture.Entry, opts Options) []Score {
	if len(entries) == 0 {
		return nil
	}

	now := time.Now()
	scores := make([]Score, len(entries))

	for i, e := range entries {
		scores[i] = Score{
			Entry: e,
			Value: compute(e, now, opts),
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Value > scores[j].Value
	})

	return scores
}

// Top returns the top n entries by prestige score.
func Top(entries []capture.Entry, n int, opts Options) []Score {
	ranked := Rank(entries, opts)
	if n > len(ranked) {
		n = len(ranked)
	}
	return ranked[:n]
}

func compute(e capture.Entry, now time.Time, opts Options) float64 {
	return opts.RecencyWeight*recencyScore(e, now) +
		opts.LatencyWeight*latencyScore(e, opts.LatencyThresholdMS) +
		opts.ErrorWeight*errorScore(e)
}

func recencyScore(e capture.Entry, now time.Time) float64 {
	if e.Timestamp.IsZero() {
		return 0
	}
	age := now.Sub(e.Timestamp).Seconds()
	// decay: score approaches 0 as age grows; full score within ~10s
	return math.Exp(-age / 60)
}

func latencyScore(e capture.Entry, thresholdMS float64) float64 {
	if thresholdMS <= 0 {
		return 0
	}
	ratio := float64(e.LatencyMS) / thresholdMS
	if ratio >= 1 {
		return 1
	}
	return ratio
}

func errorScore(e capture.Entry) float64 {
	if e.StatusCode != 0 {
		return 1
	}
	return 0
}
