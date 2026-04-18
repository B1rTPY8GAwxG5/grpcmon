// Package compare provides utilities for comparing two capture stores
// or snapshots side-by-side, producing a structured report.
package compare

import (
	"fmt"
	"io"

	"grpcmon/internal/capture"
	"grpcmon/internal/diff"
)

// Report holds the result of comparing two sets of entries.
type Report struct {
	Results []diff.Result
	MatchCount    int
	MismatchCount int
}

// Stores compares entries from two capture stores by matching on Method.
// Entries in baseline with no counterpart in candidate are skipped.
func Stores(baseline, candidate *capture.Store) Report {
	baseEntries := baseline.List()
	candEntries := candidate.List()

	candByMethod := make(map[string]capture.Entry, len(candEntries))
	for _, e := range candEntries {
		candByMethod[e.Method] = e
	}

	var report Report
	for _, base := range baseEntries {
		cand, ok := candByMethod[base.Method]
		if !ok {
			continue
		}
		r := diff.Compare(base, cand)
		report.Results = append(report.Results, r)
		if r.Match {
			report.MatchCount++
		} else {
			report.MismatchCount++
		}
	}
	return report
}

// Fprint writes a human-readable summary of the report to w.
func Fprint(w io.Writer, r Report) {
	total := r.MatchCount + r.MismatchCount
	fmt.Fprintf(w, "Compared %d entries: %d match, %d mismatch\n", total, r.MatchCount, r.MismatchCount)
	for _, res := range r.Results {
		fmt.Fprintln(w, res.String())
	}
}
