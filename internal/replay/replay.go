package replay

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/grpcmon/internal/capture"
)

// Result holds the outcome of replaying a single entry.
type Result struct {
	EntryID   string
	Method    string
	Duration  time.Duration
	StatusCode int
	Err       error
}

// Options controls replay behaviour.
type Options struct {
	// DelayBetween adds a pause between replayed requests.
	DelayBetween time.Duration
	// TimeoutPerRequest caps each individual request.
	TimeoutPerRequest time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		DelayBetween:      0,
		TimeoutPerRequest: 10 * time.Second,
	}
}

// Replayer sends captured entries to a gRPC target.
type Replayer struct {
	conn *grpc.ClientConn
	opts Options
}

// New creates a Replayer connected to addr.
func New(addr string, opts Options) (*Replayer, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) //nolint:staticcheck
	if err != nil {
		return nil, fmt.Errorf("replay: dial %s: %w", addr, err)
	}
	return &Replayer{conn: conn, opts: opts}, nil
}

// Close tears down the underlying connection.
func (r *Replayer) Close() error {
	return r.conn.Close()
}

// ReplayEntry replays a single captured entry and returns its result.
func (r *Replayer) ReplayEntry(ctx context.Context, e capture.Entry) Result {
	result := Result{
		EntryID: e.ID,
		Method:  e.Method,
	}

	if r.opts.TimeoutPerRequest > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.opts.TimeoutPerRequest)
		defer cancel()
	}

	if len(e.RequestMeta) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(e.RequestMeta))
	}

	start := time.Now()
	var reply interface{}
	err := r.conn.Invoke(ctx, e.Method, e.Request, &reply)
	result.Duration = time.Since(start)
	result.Err = err

	return result
}

// ReplayAll replays a slice of entries, streaming results into the returned channel.
func (r *Replayer) ReplayAll(ctx context.Context, entries []capture.Entry) <-chan Result {
	ch := make(chan Result, len(entries))
	go func() {
		defer close(ch)
		for _, e := range entries {
			select {
			case <-ctx.Done():
				return
			default:
			}
			ch <- r.ReplayEntry(ctx, e)
			if r.opts.DelayBetween > 0 {
				time.Sleep(r.opts.DelayBetween)
			}
		}
	}()
	return ch
}
