// Package diff provides comparison utilities for grpcmon capture entries.
//
// It is intended for use in development workflows where a developer wants to
// replay captured gRPC traffic and verify that responses have not changed —
// for example, after a server-side refactor.
//
// Basic usage:
//
//	result := diff.Compare(original, replayed)
//	if !result.Match {
//		fmt.Println(result)
//	}
//
// For bulk comparison of two slices of entries captured at different times:
//
//	results := diff.CompareAll(baseline, current)
//	for i, r := range results {
//		if !r.Match {
//			fmt.Printf("entry %d: %s\n", i, r)
//		}
//	}
package diff
