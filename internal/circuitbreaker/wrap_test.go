package circuitbreaker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/grpcmon/internal/capture"
	"github.com/user/grpcmon/internal/circuitbreaker"
)

func makeEntry(method string) capture.Entry {
	return capture.Entry{Method: method}
}

func TestWrap_PassesThroughWhenClosed(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.Options{MaxFailures: 3, Cooldown: time.Minute})
	called := false
	fn := circuitbreaker.Wrap(b, func(_ context.Context, _ capture.Entry) error {
		called = true
		return nil
	})

	if err := fn(context.Background(), makeEntry("/svc/Method")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected inner function to be called")
	}
}

func TestWrap_RecordsFailure(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.Options{MaxFailures: 2, Cooldown: time.Minute})
	inner := errors.New("rpc error")
	fn := circuitbreaker.Wrap(b, func(_ context.Context, _ capture.Entry) error {
		return inner
	})

	_ = fn(context.Background(), makeEntry("/svc/A"))
	_ = fn(context.Background(), makeEntry("/svc/A"))

	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected circuit to open after failures")
	}
}

func TestWrap_RejectsWhenOpen(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.Options{MaxFailures: 1, Cooldown: time.Minute})
	b.RecordFailure() // force open

	var called bool
	fn := circuitbreaker.Wrap(b, func(_ context.Context, _ capture.Entry) error {
		called = true
		return nil
	})

	err := fn(context.Background(), makeEntry("/svc/B"))
	if !errors.Is(err, circuitbreaker.ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
	if called {
		t.Fatal("inner function must not be called when circuit is open")
	}
}

func TestWrap_RecordsSuccessAfterHalfOpen(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.Options{MaxFailures: 1, Cooldown: time.Millisecond})
	b.RecordFailure()
	time.Sleep(5 * time.Millisecond)

	fn := circuitbreaker.Wrap(b, func(_ context.Context, _ capture.Entry) error {
		return nil
	})

	if err := fn(context.Background(), makeEntry("/svc/C")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected closed, got %v", b.State())
	}
}
