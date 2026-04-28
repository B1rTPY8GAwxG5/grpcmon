package decay_test

import (
	"context"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/decay"
)

func noopReplayer(called *bool) decay.Replayer {
	return decay.ReplayerFunc(func(_ context.Context, e capture.Entry) (capture.Entry, error) {
		*called = true
		e.Response = "replayed"
		return e, nil
	})
}

func TestWrap_AllowsRecentEntry(t *testing.T) {
	var called bool
	s := decay.New(decay.DefaultOptions())
	w := decay.Wrap(noopReplayer(&called), s, 0.5)

	e := makeEntry("/svc/Fresh", 0)
	out, err := w.ReplayEntry(context.Background(), e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected replayer to be called for recent entry")
	}
	if out.Response != "replayed" {
		t.Fatalf("expected replayed response, got %v", out.Response)
	}
}

func TestWrap_SkipsStaleEntry(t *testing.T) {
	var called bool
	opts := decay.Options{HalfLife: time.Millisecond}
	s := decay.New(opts)
	w := decay.Wrap(noopReplayer(&called), s, 0.5)

	e := makeEntry("/svc/Stale", 10*time.Second)
	out, err := w.ReplayEntry(context.Background(), e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected replayer to be skipped for stale entry")
	}
	// original entry returned unchanged
	if out.Method != "/svc/Stale" {
		t.Fatalf("expected original entry returned, got %v", out.Method)
	}
}

func TestWrap_ZeroThreshold_AlwaysReplays(t *testing.T) {
	var called bool
	opts := decay.Options{HalfLife: time.Millisecond}
	s := decay.New(opts)
	w := decay.Wrap(noopReplayer(&called), s, 0)

	e := makeEntry("/svc/Ancient", 24*time.Hour)
	_, err := w.ReplayEntry(context.Background(), e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected replayer called with zero threshold")
	}
}
