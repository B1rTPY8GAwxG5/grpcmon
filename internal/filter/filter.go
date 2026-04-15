package filter

import (
	"strings"

	"github.com/grpcmon/internal/capture"
)

// Criteria holds the filtering parameters for captured RPC entries.
type Criteria struct {
	Method     string
	StatusCode string
	MinLatency int64 // milliseconds
	MaxLatency int64 // milliseconds, 0 means no upper bound
}

// Match reports whether the given entry satisfies all non-zero criteria.
func Match(entry capture.Entry, c Criteria) bool {
	if c.Method != "" && !strings.Contains(entry.Method, c.Method) {
		return false
	}

	if c.StatusCode != "" && !strings.EqualFold(entry.StatusCode, c.StatusCode) {
		return false
	}

	latencyMs := entry.Duration.Milliseconds()

	if c.MinLatency > 0 && latencyMs < c.MinLatency {
		return false
	}

	if c.MaxLatency > 0 && latencyMs > c.MaxLatency {
		return false
	}

	return true
}

// Apply returns the subset of entries from the store that satisfy criteria.
func Apply(entries []capture.Entry, c Criteria) []capture.Entry {
	result := make([]capture.Entry, 0, len(entries))
	for _, e := range entries {
		if Match(e, c) {
			result = append(result, e)
		}
	}
	return result
}
