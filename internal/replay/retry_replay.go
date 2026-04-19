package replay

import (
	"context"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/retry"
)

// ReplayWithRetry replays a single entry using the provided Replayer, applying
// the given retry policy on transient failures.
func ReplayWithRetry(ctx context.Context, r *Replayer, e capture.Entry, p retry.Policy) (Result, error) {
	var res Result
	attempts, err := retry.Do(ctx, p, func() error {
		var rerr error
		res, rerr = r.ReplayEntry(ctx, e)
		return rerr
	})
	res.Attempts = attempts
	return res, err
}

// Result extends the base diff result with attempt metadata.
type Result struct {
	Entry    capture.Entry
	Err      error
	Attempts int
}
