package cooldown

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/example/grpcmon/internal/capture"
)

func makeEntry(method string) capture.Entry {
	return capture.Entry{Method: method}
}

func noopReplayer(_ context.Context, e capture.Entry) (capture.Entry, error) {
	return e, nil
}

func TestWrap_AllowsFirstCall(t *testing.T) {
	c := New(500 * time.Millisecond)
	r := Wrap(noopReplayer, c)
	_, err := r(context.Background(), makeEntry("/svc/M"))
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
}

func TestWrap_RejectsWithinCooldown(t *testing.T) {
	c := New(500 * time.Millisecond)
	r := Wrap(noopReplayer, c)
	r(context.Background(), makeEntry("/svc/M")) //nolint:errcheck
	_, err := r(context.Background(), makeEntry("/svc/M"))
	if err == nil {
		t.Fatal("expected error on second call within cooldown")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T", err)
	}
	if st.Code() != codes.ResourceExhausted {
		t.Fatalf("expected ResourceExhausted, got %v", st.Code())
	}
}

func TestWrap_AllowsAfterCooldown(t *testing.T) {
	c := New(20 * time.Millisecond)
	r := Wrap(noopReplayer, c)
	r(context.Background(), makeEntry("/svc/M")) //nolint:errcheck
	time.Sleep(30 * time.Millisecond)
	_, err := r(context.Background(), makeEntry("/svc/M"))
	if err != nil {
		t.Fatalf("expected nil error after cooldown elapsed: %v", err)
	}
}

func TestWrap_PropagatesReplayerError(t *testing.T) {
	c := New(10 * time.Millisecond)
	sentinel := errors.New("downstream error")
	errReplayer := func(_ context.Context, e capture.Entry) (capture.Entry, error) {
		return capture.Entry{}, sentinel
	}
	r := Wrap(errReplayer, c)
	_, err := r(context.Background(), makeEntry("/svc/M"))
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestWrap_NilCooldown_UsesDefault(t *testing.T) {
	r := Wrap(noopReplayer, nil)
	_, err := r(context.Background(), makeEntry("/svc/M"))
	if err != nil {
		t.Fatalf("unexpected error with nil cooldown: %v", err)
	}
}
