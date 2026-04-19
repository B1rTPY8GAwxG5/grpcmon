package metric

import (
	"fmt"
	"io"
	"sort"
)

// MethodSummary holds a rolled-up view across all windows for a method.
type MethodSummary struct {
	Method     string
	Total      int
	Errors     int
	ErrorRate  float64
	AvgLatency float64
}

// Summarise computes a MethodSummary for each tracked method.
func (t *Tracker) Summarise() []MethodSummary {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make([]MethodSummary, 0, len(t.buckets))
	for method, windows := range t.buckets {
		var total, errors int
		var latSum float64
		var latCount int
		for _, w := range windows {
			total += w.Total
			errors += w.Errors
			for _, l := range w.Latency {
				latSum += l
				latCount++
			}
		}
		var errRate, avgLat float64
		if total > 0 {
			errRate = float64(errors) / float64(total)
		}
		if latCount > 0 {
			avgLat = latSum / float64(latCount)
		}
		out = append(out, MethodSummary{
			Method:     method,
			Total:      total,
			Errors:     errors,
			ErrorRate:  errRate,
			AvgLatency: avgLat,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Method < out[j].Method })
	return out
}

// Fprint writes a human-readable metric report to w.
func Fprint(w io.Writer, summaries []MethodSummary) {
	fmt.Fprintf(w, "%-40s %8s %8s %10s %12s\n", "METHOD", "TOTAL", "ERRORS", "ERR RATE", "AVG LAT ms")
	for _, s := range summaries {
		fmt.Fprintf(w, "%-40s %8d %8d %9.1f%% %12.1f\n",
			s.Method, s.Total, s.Errors, s.ErrorRate*100, s.AvgLatency)
	}
}
