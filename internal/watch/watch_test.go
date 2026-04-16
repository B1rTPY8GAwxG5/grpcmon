package watch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/watch"
)

func TestWatcher_NotifiesOnNewEntries(t *testing.T) {
	store := capture.NewStore(100)

	var mu sync.Mutex
	var received []capture.Entry

	w := watch.New(store, 20*time.Millisecond, func(entries []capture.Entry) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, entries...)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.Run(ctx)

	store.Add(capture.Entry{ID: "1", Method: "/svc/A"})
	store.Add(capture.Entry{ID: "2", Method: "/svc/B"})

	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) < 2 {
		t.Fatalf("expected at least 2 entries, got %d", len(received))
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	store := capture.NewStore(100)
	calls := 0

	w := watch.New(store, 10*time.Millisecond, func(_ []capture.Entry) {
		calls++
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancel")
	}
}

func TestWatcher_NoCallWhenNoNewEntries(t *testing.T) {
	store := capture.NewStore(100)
	store.Add(capture.Entry{ID: "x", Method: "/svc/X"})

	calls := 0
	w := watch.New(store, 20*time.Millisecond, func(_ []capture.Entry) {
		calls++
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Run(ctx)

	// consume the initial entry
	time.Sleep(50 * time.Millisecond)
	snap := calls

	// no new entries added — calls should not grow
	time.Sleep(60 * time.Millisecond)
	if calls > snap+1 {
		t.Fatalf("handler called unexpectedly: %d times after snapshot", calls-snap)
	}
}
