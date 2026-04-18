// Package alert evaluates captured gRPC entries against user-defined rules
// and produces alerts when error rates or latencies exceed configured
// thresholds. Rules can be scoped to a specific gRPC method or applied
// globally across all captured traffic.
package alert
