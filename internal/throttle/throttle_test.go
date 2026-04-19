package throttle_test

import (
	"context"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/throttle"
)

func makeEntries(offsets ...time.Duration) []capture.Entry {
	base := time.Now()
	entries := make([]capture.Entry, len(offsets))
	for i, d := range offsets {
		entries[i] = capture.Entry{Timestamp: base.Add(d)}
	}
	return entries
}

func TestRun_CallsReplayerForEachEntry(t *testing.T) {
	entries := makeEntries(0, 10*time.Millisecond, 20*time.Millisecond)
	called := 0
	err := throttle.Run(context.Background(), entries, func(_ context.Context, _ capture.Entry) error {
		called++
		return nil
	}, throttle.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 3 {
		t.Fatalf("expected 3 calls, got %d", called)
	}
}

func TestRun_RespectsContextCancel(t *testing.T) {
	entries := makeEntries(0, 500*time.Millisecond, time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	called := 0
	cancel()
	err := throttle.Run(ctx, entries, func(_ context.Context, _ capture.Entry) error {
		called++
		return nil
	}, throttle.DefaultOptions())
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestRun_MaxDelayCaps(t *testing.T) {
	entries := makeEntries(0, 10*time.Second)
	opts := throttle.Options{SpeedFactor: 1.0, MaxDelay: 10 * time.Millisecond}
	start := time.Now()
	_ = throttle.Run(context.Background(), entries, func(_ context.Context, _ capture.Entry) error {
		return nil
	}, opts)
	if time.Since(start) > 200*time.Millisecond {
		t.Fatal("MaxDelay was not respected")
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := throttle.DefaultOptions()
	if opts.SpeedFactor != 1.0 {
		t.Errorf("expected SpeedFactor 1.0, got %v", opts.SpeedFactor)
	}
	if opts.MaxDelay != 5*time.Second {
		t.Errorf("expected MaxDelay 5s, got %v", opts.MaxDelay)
	}
}
