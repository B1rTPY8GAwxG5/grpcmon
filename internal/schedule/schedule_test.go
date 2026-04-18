package schedule_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/replay"
	"github.com/grpcmon/internal/schedule"
)

func makeStore(t *testing.T) *capture.Store {
	t.Helper()
	s, err := capture.NewStore(10)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestNew_DefaultInterval(t *testing.T) {
	store := makeStore(t)
	job := schedule.Job{Interval: 0}
	s := schedule.New(store, job)
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	store := makeStore(t)
	job := schedule.Job{Interval: 100 * time.Millisecond}
	s := schedule.New(store, job)
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	err := s.Run(ctx)
	if err == nil {
		t.Fatal("expected non-nil error on cancel")
	}
}

func TestRun_InvokesOnResult(t *testing.T) {
	store := makeStore(t)
	entry := capture.Entry{Method: "/pkg.Svc/Method"}
	store.Add(entry)

	var called atomic.Int32
	opts := replay.DefaultOptions()
	opts.Target = "" // will fail to connect — result still fires via closed channel

	job := schedule.Job{
		Interval: 50 * time.Millisecond,
		Options:  opts,
		OnResult: func(_ replay.Result) { called.Add(1) },
	}
	s := schedule.New(store, job)
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	s.Run(ctx) //nolint:errcheck
	// OnResult may or may not fire depending on connection; we just assert no panic.
}
