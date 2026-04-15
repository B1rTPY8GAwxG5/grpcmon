package replay

import "time"

// Options configures the behaviour of a Replayer.
type Options struct {
	// Timeout is the per-call deadline applied when replaying an entry.
	Timeout time.Duration

	// Concurrency is the maximum number of parallel replay goroutines used by
	// ReplayAll. A value of 0 or less is treated as 1 (sequential).
	Concurrency int

	// DelayBetweenCalls is an optional pause inserted between successive calls
	// in ReplayAll when Concurrency == 1.
	DelayBetweenCalls time.Duration
}

// DefaultOptions returns an Options struct populated with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Timeout:           5 * time.Second,
		Concurrency:       1,
		DelayBetweenCalls: 0,
	}
}
