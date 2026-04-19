package retry

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Policy defines retry behaviour.
type Policy struct {
	MaxAttempts int
	Backoff     time.Duration
	RetryOn     []codes.Code
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Backoff:     200 * time.Millisecond,
		RetryOn:     []codes.Code{codes.Unavailable, codes.DeadlineExceeded},
	}
}

// Do executes fn up to p.MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or the error code is not retryable.
func Do(ctx context.Context, p Policy, fn func() error) (int, error) {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	var last error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return attempt - 1, err
		}
		last = fn()
		if last == nil {
			return attempt, nil
		}
		if !isRetryable(last, p.RetryOn) {
			return attempt, last
		}
		if attempt < p.MaxAttempts {
			select {
			case <-ctx.Done():
				return attempt, ctx.Err()
			case <-time.After(p.Backoff):
			}
		}
	}
	return p.MaxAttempts, last
}

func isRetryable(err error, codes []codes.Code) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	for _, c := range codes {
		if st.Code() == c {
			return true
		}
	}
	return false
}
