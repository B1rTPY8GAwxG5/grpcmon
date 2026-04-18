// Package compare provides side-by-side comparison of two capture stores
// or snapshots, using the diff package to detect mismatches between
// baseline and candidate gRPC traffic recordings.
//
// Typical usage:
//
//	report := compare.Stores(baselineStore, candidateStore)
//	compare.Fprint(os.Stdout, report)
package compare
