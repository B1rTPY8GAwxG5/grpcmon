// Package debounce coalesces rapid bursts of capture.Entry values into
// batched handler invocations.
//
// During high-frequency gRPC traffic a watcher or TUI refresh loop may be
// triggered hundreds of times per second. Debounce delays notification until
// a configurable period of inactivity has passed, then delivers all
// accumulated entries in a single call.
//
// Basic usage:
//
//	d := debounce.New(100*time.Millisecond, func(entries []capture.Entry) {
//		// process batch
//	})
//	d.Run(ctx)          // flush on shutdown
//	d.Add(entry)        // queue an entry
//
// Flush can be called at any time to drain pending entries immediately,
// which is useful in tests or on graceful shutdown.
package debounce
