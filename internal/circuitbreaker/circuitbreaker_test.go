package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/grpcmon/internal/circuitbreaker"
)

func TestDefaultOptions_Values(t *testing.T) {
	opts := circuitbreaker.DefaultOptions()
	if opts.MaxFailures != 5 {
		t.Fatalf("expected MaxFailures=5, got %d", opts.MaxFailures)
	}
	if opts.Cooldown != 10*time.Second {
		t.Fatalf("expected Cooldown=10s, got %v", opts.Cooldown)
	}
}

func TestBreaker_InitiallyClosed(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.DefaultOptions())
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected circuit to be closed initially")
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil error when closed, got %v", err)
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	opts := circuitbreaker.Options{MaxFailures: 3, Cooldown: time.Minute}
	b := circuitbreaker.New(opts)

	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}

	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected circuit to be open after threshold")
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_ClosesAfterSuccess(t *testing.T) {
	opts := circuitbreaker.Options{MaxFailures: 1, Cooldown: time.Millisecond}
	b := circuitbreaker.New(opts)
	b.RecordFailure()

	time.Sleep(5 * time.Millisecond)

	if err := b.Allow(); err != nil {
		t.Fatalf("expected half-open allow, got %v", err)
	}
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected circuit closed after success in half-open")
	}
}

func TestBreaker_ReopensFromHalfOpen(t *testing.T) {
	opts := circuitbreaker.Options{MaxFailures: 1, Cooldown: time.Millisecond}
	b := circuitbreaker.New(opts)
	b.RecordFailure()

	time.Sleep(5 * time.Millisecond)
	_ = b.Allow() // transitions to half-open
	b.RecordFailure()

	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected circuit to reopen after failure in half-open")
	}
}

func TestBreaker_SuccessResetsFailureCount(t *testing.T) {
	opts := circuitbreaker.Options{MaxFailures: 3, Cooldown: time.Minute}
	b := circuitbreaker.New(opts)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()

	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected circuit to remain closed after reset")
	}
}
