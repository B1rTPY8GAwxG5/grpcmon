// Package summary provides a concise text summary of captured gRPC traffic.
package summary

import (
	"fmt"
	"io"
	"time"

	"grpcmon/internal/stats"
)

// Report holds a formatted summary derived from a Stats value.
type Report struct {
	Total      int
	Successful int
	Failed     int
	ErrorRate  float64
	P50        time.Duration
	P99        time.Duration
	TopMethod  string
}

// FromStats builds a Report from a stats.Stats value.
func FromStats(s stats.Stats) Report {
	top := ""
	if len(s.TopMethods) > 0 {
		top = s.TopMethods[0].Method
	}
	return Report{
		Total:      s.Total,
		Successful: s.Successful,
		Failed:     s.Failed,
		ErrorRate:  s.ErrorRate,
		P50:        s.P50,
		P99:        s.P99,
		TopMethod:  top,
	}
}

// Fprint writes a human-readable summary to w.
func Fprint(w io.Writer, r Report) {
	fmt.Fprintf(w, "Total requests : %d\n", r.Total)
	fmt.Fprintf(w, "Successful     : %d\n", r.Successful)
	fmt.Fprintf(w, "Failed         : %d\n", r.Failed)
	fmt.Fprintf(w, "Error rate     : %.1f%%\n", r.ErrorRate*100)
	fmt.Fprintf(w, "Latency p50    : %s\n", r.P50.Round(time.Millisecond))
	fmt.Fprintf(w, "Latency p99    : %s\n", r.P99.Round(time.Millisecond))
	if r.TopMethod != "" {
		fmt.Fprintf(w, "Top method     : %s\n", r.TopMethod)
	}
}
