package jitter_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/jitter"
)

func noopReplayer(_ context.Context, _ capture.Entry) error { return nil }

func TestDefaultOptions_Values(t *testing.T) {
	opts := jitter.DefaultOptions()
	if opts.MinDelay != 0 {
		t.Fatalf("expected MinDelay 0, got %v", opts.MinDelay)
	}
	if opts.MaxDelay != 250*time.Millisecond {
		t.Fatalf("expected MaxDelay 250ms, got %v", opts.MaxDelay)
	}
}

func TestWrap_ZeroMaxDelay_NoDelay(t *testing.T) {
	opts := jitter.Options{MinDelay: 0, MaxDelay: 0}
	r := jitter.Wrap(noopReplayer, opts)

	start := time.Now()
	if err := r(context.Background(), capture.Entry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Fatalf("expected near-zero delay, got %v", elapsed)
	}
}

func TestWrap_AddsDelay(t *testing.T) {
	opts := jitter.Options{
		MinDelay: 20 * time.Millisecond,
		MaxDelay: 40 * time.Millisecond,
		Seed:     42,
	}
	r := jitter.Wrap(noopReplayer, opts)

	start := time.Now()
	if err := r(context.Background(), capture.Entry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < 20*time.Millisecond {
		t.Fatalf("expected at least 20ms delay, got %v", elapsed)
	}
}

func TestWrap_RespectsContextCancel(t *testing.T) {
	opts := jitter.Options{
		MinDelay: 500 * time.Millisecond,
		MaxDelay: 1 * time.Second,
		Seed:     1,
	}
	r := jitter.Wrap(noopReplayer, opts)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r(ctx, capture.Entry{})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestWrap_DelegatesEntry(t *testing.T) {
	opts := jitter.Options{MinDelay: 0, MaxDelay: 1 * time.Millisecond, Seed: 7}

	want := "test.method"
	var got string

	capturing := func(_ context.Context, e capture.Entry) error {
		got = e.Method
		return nil
	}

	r := jitter.Wrap(capturing, opts)
	e := capture.Entry{Method: want}
	if err := r(context.Background(), e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("expected method %q, got %q", want, got)
	}
}
