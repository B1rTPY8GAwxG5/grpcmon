// Package throttle provides replay throttling based on original request timing.
package throttle

import (
	"context"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Replayer is a function that replays a single entry.
type Replayer func(ctx context.Context, entry capture.Entry) error

// Options configures throttle behaviour.
type Options struct {
	// SpeedFactor scales the delay between replayed requests.
	// 1.0 = real time, 2.0 = double speed, 0.5 = half speed.
	SpeedFactor float64
	// MaxDelay caps the inter-request delay regardless of SpeedFactor.
	MaxDelay time.Duration
}

// DefaultOptions returns sensible throttle defaults.
func DefaultOptions() Options {
	return Options{
		SpeedFactor: 1.0,
		MaxDelay:    5 * time.Second,
	}
}

// Run replays entries in timestamp order, preserving relative timing.
func Run(ctx context.Context, entries []capture.Entry, fn Replayer, opts Options) error {
	if opts.SpeedFactor <= 0 {
		opts.SpeedFactor = 1.0
	}
	if opts.MaxDelay <= 0 {
		opts.MaxDelay = DefaultOptions().MaxDelay
	}

	for i, entry := range entries {
		if i > 0 {
			prev := entries[i-1]
			gap := entry.Timestamp.Sub(prev.Timestamp)
			if gap > 0 {
				delay := time.Duration(float64(gap) / opts.SpeedFactor)
				if delay > opts.MaxDelay {
					delay = opts.MaxDelay
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}
			}
		}
		if err := fn(ctx, entry); err != nil {
			return err
		}
	}
	return nil
}
