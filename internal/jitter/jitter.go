// Package jitter adds randomised delay to replay operations to avoid
// thundering-herd effects when replaying captured traffic.
package jitter

import (
	"context"
	"math/rand"
	"time"

	"github.com/example/grpcmon/internal/capture"
)

// Options configures jitter behaviour.
type Options struct {
	// MinDelay is the minimum additional delay applied before each replay.
	MinDelay time.Duration
	// MaxDelay is the maximum additional delay applied before each replay.
	MaxDelay time.Duration
	// Seed is used to initialise the random source. Zero means use a
	// time-based seed.
	Seed int64
}

// DefaultOptions returns sensible defaults for jitter.
func DefaultOptions() Options {
	return Options{
		MinDelay: 0,
		MaxDelay: 250 * time.Millisecond,
	}
}

// Replayer is a function that replays a single captured entry.
type Replayer func(ctx context.Context, e capture.Entry) error

// Wrap returns a Replayer that waits a random duration in [min, max] before
// delegating to next. If min == max the delay is constant.
func Wrap(next Replayer, opts Options) Replayer {
	if opts.MaxDelay <= 0 {
		return next
	}

	min := opts.MinDelay
	max := opts.MaxDelay
	if min > max {
		min, max = max, min
	}

	src := rand.NewSource(opts.Seed)
	if opts.Seed == 0 {
		src = rand.NewSource(time.Now().UnixNano())
	}
	rng := rand.New(src) //nolint:gosec

	window := int64(max - min)

	return func(ctx context.Context, e capture.Entry) error {
		delay := min
		if window > 0 {
			delay += time.Duration(rng.Int63n(window))
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
		return next(ctx, e)
	}
}
