package circuitbreaker

import (
	"context"

	"github.com/user/grpcmon/internal/capture"
)

// ReplayFunc is a function that replays a single capture entry.
type ReplayFunc func(ctx context.Context, entry capture.Entry) error

// Wrap returns a ReplayFunc that gates calls through the Breaker. If the
// circuit is open the original entry is skipped and ErrOpen is returned.
// On success the breaker is reset; on failure the failure counter advances.
func Wrap(b *Breaker, fn ReplayFunc) ReplayFunc {
	return func(ctx context.Context, entry capture.Entry) error {
		if err := b.Allow(); err != nil {
			return err
		}
		err := fn(ctx, entry)
		if err != nil {
			b.RecordFailure()
			return err
		}
		b.RecordSuccess()
		return nil
	}
}
