// Package jitter provides a Replayer wrapper that introduces a randomised
// delay before each replay call.
//
// Use jitter when replaying captured traffic against a live service to spread
// load more naturally and avoid synchronised bursts that would not occur in
// real production traffic.
//
// Example:
//
//	opts := jitter.DefaultOptions()
//	opts.MaxDelay = 100 * time.Millisecond
//	delayed := jitter.Wrap(myReplayer, opts)
//	// delayed now waits 0–100 ms before each call
package jitter
