// Package stats provides aggregation and summary statistics over captured gRPC entries.
package stats

import (
	"time"

	"github.com/grpcmon/internal/capture"
	"google.golang.org/grpc/codes"
)

// Summary holds aggregated statistics for a collection of captured entries.
type Summary struct {
	Total        int
	SuccessCount int
	ErrorCount   int
	AvgLatency   time.Duration
	MinLatency   time.Duration
	MaxLatency   time.Duration
	StatusCodes  map[codes.Code]int
	TopMethods   []MethodStat
}

// MethodStat holds call count for a single gRPC method.
type MethodStat struct {
	Method string
	Count  int
}

// Compute derives a Summary from the provided entries.
func Compute(entries []capture.Entry) Summary {
	if len(entries) == 0 {
		return Summary{StatusCodes: make(map[codes.Code]int)}
	}

	statusCodes := make(map[codes.Code]int)
	methodCounts := make(map[string]int)

	var totalLatency time.Duration
	minLat := entries[0].Latency
	maxLat := entries[0].Latency

	for _, e := range entries {
		statusCodes[e.StatusCode]++
		methodCounts[e.Method]++

		totalLatency += e.Latency
		if e.Latency < minLat {
			minLat = e.Latency
		}
		if e.Latency > maxLat {
			maxLat = e.Latency
		}

		if e.StatusCode == codes.OK {
			// counted below via statusCodes
		}
	}

	success := statusCodes[codes.OK]

	return Summary{
		Total:        len(entries),
		SuccessCount: success,
		ErrorCount:   len(entries) - success,
		AvgLatency:   totalLatency / time.Duration(len(entries)),
		MinLatency:   minLat,
		MaxLatency:   maxLat,
		StatusCodes:  statusCodes,
		TopMethods:   topN(methodCounts, 5),
	}
}

// topN returns up to n methods sorted by call count descending.
func topN(counts map[string]int, n int) []MethodStat {
	stats := make([]MethodStat, 0, len(counts))
	for m, c := range counts {
		stats = append(stats, MethodStat{Method: m, Count: c})
	}
	// simple insertion sort — entry counts are small in dev environments
	for i := 1; i < len(stats); i++ {
		for j := i; j > 0 && stats[j].Count > stats[j-1].Count; j-- {
			stats[j], stats[j-1] = stats[j-1], stats[j]
		}
	}
	if len(stats) > n {
		return stats[:n]
	}
	return stats
}
