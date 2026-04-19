package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDefaultPolicy_Values(t *testing.T) {
	p := DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected 3 attempts, got %d", p.MaxAttempts)
	}
	if p.Backoff != 200*time.Millisecond {
		t.Errorf("unexpected backoff %v", p.Backoff)
	}
}

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	p := DefaultPolicy()
	p.Backoff = 0
	attempts, err := Do(context.Background(), p, func() error { return nil })
	if err != nil || attempts != 1 {
		t.Errorf("expected 1 attempt and no error, got %d %v", attempts, err)
	}
}

func TestDo_RetriesOnRetryableCode(t *testing.T) {
	p := DefaultPolicy()
	p.Backoff = 0
	calls := 0
	_, err := Do(context.Background(), p, func() error {
		calls++
		return status.Error(codes.Unavailable, "down")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != p.MaxAttempts {
		t.Errorf("expected %d calls, got %d", p.MaxAttempts, calls)
	}
}

func TestDo_StopsOnNonRetryableCode(t *testing.T) {
	p := DefaultPolicy()
	p.Backoff = 0
	calls := 0
	_, err := Do(context.Background(), p, func() error {
		calls++
		return status.Error(codes.NotFound, "missing")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_StopsOnNonGRPCError(t *testing.T) {
	p := DefaultPolicy()
	p.Backoff = 0
	calls := 0
	_, err := Do(context.Background(), p, func() error {
		calls++
		return errors.New("plain error")
	})
	if err == nil || calls != 1 {
		t.Errorf("expected 1 call and error, got %d %v", calls, err)
	}
}

func TestDo_RespectsContextCancel(t *testing.T) {
	p := DefaultPolicy()
	p.Backoff = time.Second
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Do(ctx, p, func() error {
		return status.Error(codes.Unavailable, "down")
	})
	if err == nil {
		t.Fatal("expected context error")
	}
}
