package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/timeout"
)

func makeEntry(method string) capture.Entry {
	return capture.Entry{Method: method}
}

func TestNew_DefaultsToFiveSeconds(t *testing.T) {
	m := timeout.New(0)
	if got := m.Get("any"); got != 5*time.Second {
		t.Fatalf("expected 5s default, got %v", got)
	}
}

func TestGet_ReturnsDefault_WhenMethodNotRegistered(t *testing.T) {
	m := timeout.New(2 * time.Second)
	if got := m.Get("/svc/Method"); got != 2*time.Second {
		t.Fatalf("expected 2s, got %v", got)
	}
}

func TestSet_And_Get_MethodSpecific(t *testing.T) {
	m := timeout.New(2 * time.Second)
	m.Set("/svc/Slow", 10*time.Second)
	if got := m.Get("/svc/Slow"); got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
	// Other methods still use the default.
	if got := m.Get("/svc/Fast"); got != 2*time.Second {
		t.Fatalf("expected default 2s for unregistered method, got %v", got)
	}
}

func TestWrap_CompletesWithinTimeout(t *testing.T) {
	m := timeout.New(500 * time.Millisecond)
	wrapped := m.Wrap(func(ctx context.Context, e capture.Entry) error {
		return nil // instant
	})
	if err := wrapped(context.Background(), makeEntry("/svc/OK")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWrap_ReturnsDeadlineExceeded_WhenSlowReplayer(t *testing.T) {
	m := timeout.New(20 * time.Millisecond)
	wrapped := m.Wrap(func(ctx context.Context, e capture.Entry) error {
		select {
		case <-time.After(200 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	err := wrapped(context.Background(), makeEntry("/svc/Slow"))
	if !errors.Is(err, timeout.ErrDeadlineExceeded) {
		t.Fatalf("expected ErrDeadlineExceeded, got %v", err)
	}
}

func TestWrap_UsesMethodSpecificTimeout(t *testing.T) {
	m := timeout.New(10 * time.Millisecond)
	m.Set("/svc/Generous", 500*time.Millisecond)

	wrapped := m.Wrap(func(ctx context.Context, e capture.Entry) error {
		time.Sleep(30 * time.Millisecond)
		return nil
	})

	// Should succeed because method-specific timeout is generous.
	if err := wrapped(context.Background(), makeEntry("/svc/Generous")); err != nil {
		t.Fatalf("unexpected error with generous timeout: %v", err)
	}
}

func TestWrap_PropagatesReplayerError(t *testing.T) {
	sentinel := errors.New("replayer error")
	m := timeout.New(500 * time.Millisecond)
	wrapped := m.Wrap(func(_ context.Context, _ capture.Entry) error {
		return sentinel
	})
	if err := wrapped(context.Background(), makeEntry("/svc/Err")); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
