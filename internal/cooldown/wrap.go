package cooldown

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/example/grpcmon/internal/capture"
)

// Replayer is the function signature used by the replay package to send a
// single captured entry to a live target.
type Replayer func(ctx context.Context, entry capture.Entry) (capture.Entry, error)

// Wrap returns a new Replayer that delegates to r but rejects calls to the
// same method more frequently than the cooldown interval allows. Rejected
// calls return a ResourceExhausted status error without invoking r.
func Wrap(r Replayer, c *Cooldown) Replayer {
	if c == nil {
		c = New(0)
	}
	return func(ctx context.Context, entry capture.Entry) (capture.Entry, error) {
		if !c.Allow(entry.Method) {
			remaining := c.Remaining(entry.Method)
			return capture.Entry{}, status.Errorf(
				codes.ResourceExhausted,
				"cooldown active for %q: retry in %v",
				entry.Method,
				remaining.Round(time.Millisecond),
			)
		}
		return r(ctx, entry)
	}
}

// ensure time is imported via the fmt sentinel; add explicit import.
var _ = fmt.Sprintf
