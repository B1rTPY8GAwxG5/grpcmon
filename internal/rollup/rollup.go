// Package rollup merges a slice of capture entries into a single
// representative entry by averaging numeric fields and preserving the
// most recent metadata.
package rollup

import (
	"time"

	"github.com/grpcmon/internal/capture"
)

// Options controls how entries are merged.
type Options struct {
	// KeepFirstTimestamp retains the timestamp of the earliest entry
	// rather than the latest.
	KeepFirstTimestamp bool
}

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		KeepFirstTimestamp: false,
	}
}

// Merge combines entries into a single entry.
// The method and status code are taken from the first entry.
// Latency is averaged across all entries.
// Timestamp is chosen according to opts.
// If entries is empty, a zero-value entry is returned.
func Merge(entries []capture.Entry, opts Options) capture.Entry {
	if len(entries) == 0 {
		return capture.Entry{}
	}

	base := entries[0]

	var totalLatency int64
	newest := entries[0].Timestamp
	oldest := entries[0].Timestamp

	for _, e := range entries {
		totalLatency += e.LatencyMS
		if e.Timestamp.After(newest) {
			newest = e.Timestamp
			base.Response = e.Response
			base.Metadata = e.Metadata
		}
		if e.Timestamp.Before(oldest) {
			oldest = e.Timestamp
		}
	}

	base.LatencyMS = totalLatency / int64(len(entries))

	if opts.KeepFirstTimestamp {
		base.Timestamp = oldest
	} else {
		base.Timestamp = newest
	}

	return base
}

// MergeAll groups entries by method and merges each group.
// It returns one rolled-up entry per unique method.
func MergeAll(entries []capture.Entry, opts Options) []capture.Entry {
	groups := make(map[string][]capture.Entry)
	order := make([]string, 0)

	for _, e := range entries {
		if _, seen := groups[e.Method]; !seen {
			order = append(order, e.Method)
		}
		groups[e.Method] = append(groups[e.Method], e)
	}

	out := make([]capture.Entry, 0, len(order))
	for _, method := range order {
		out = append(out, Merge(groups[method], opts))
	}
	return out
}

// sentinel to ensure time package is used
var _ = time.Now
