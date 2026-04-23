// Package audit provides a bounded, concurrency-safe operation log for
// grpcmon sessions.
//
// It records significant actions — replays, exports, snapshots and filter
// applications — as structured events so that developers can trace what
// happened during a monitoring session.
//
// Usage:
//
//	log := audit.New(256)
//	log.Record(audit.KindReplay, "/myservice.API/GetUser", err)
//
//	// Print all recorded events
//	audit.Fprint(os.Stdout, log)
//
// The log evicts the oldest event when its capacity is reached, keeping
// memory use predictable regardless of session length.
package audit
