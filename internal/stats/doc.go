// Package stats computes aggregated summary statistics over a slice of
// captured gRPC entries.
//
// Usage:
//
//	entries := store.List()
//	summary := stats.Compute(entries)
//	fmt.Printf("Total: %d  Success: %d  Errors: %d\n",
//		summary.Total, summary.SuccessCount, summary.ErrorCount)
//	fmt.Printf("Avg latency: %v  Min: %v  Max: %v\n",
//		summary.AvgLatency, summary.MinLatency, summary.MaxLatency)
//
// The package is intentionally stateless — call Compute with whatever
// snapshot of entries you need and discard the result when done.
package stats
