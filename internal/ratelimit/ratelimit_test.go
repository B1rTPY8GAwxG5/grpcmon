package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/grpcmon/internal/ratelimit"
)

func TestNew_AllowsTokenAcquisition(t *testing.T) {
	l := ratelimit.New(10)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected token acquisition, got: %v", err)
	}
}

func TestWait_RespectsContextCancel(t *testing.T) {
	// rps=1 with full bucket drained first so next Wait must block
	l := ratelimit.New(1)
	defer l.Stop()

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	defer cancel1()
	// drain the first token
	_ = l.Wait(ctx1)

	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()

	err := l.Wait(ctx2)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestNew_ZeroRPSDefaultsToOne(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected token with defaulted rps=1, got: %v", err)
	}
}

func TestWait_MultipleTokens(t *testing.T) {
	l := ratelimit.New(50)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
	}
}
