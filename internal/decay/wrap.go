package decay

import (
	"context"

	"github.com/grpcmon/internal/capture"
)

// Replayer is the interface satisfied by replay.Client and compatible types.
type Replayer interface {
	ReplayEntry(ctx context.Context, e capture.Entry) (capture.Entry, error)
}

// ReplayerFunc is a function adapter for Replayer.
type ReplayerFunc func(ctx context.Context, e capture.Entry) (capture.Entry, error)

// ReplayEntry implements Replayer.
func (f ReplayerFunc) ReplayEntry(ctx context.Context, e capture.Entry) (capture.Entry, error) {
	return f(ctx, e)
}

// Wrap returns a Replayer that skips entries whose decay score is below
// threshold, returning the original entry unchanged in that case.
func Wrap(r Replayer, s *Scorer, threshold float64) Replayer {
	return ReplayerFunc(func(ctx context.Context, e capture.Entry) (capture.Entry, error) {
		if s.Score(e) < threshold {
			return e, nil
		}
		return r.ReplayEntry(ctx, e)
	})
}
