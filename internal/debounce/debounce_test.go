package debounce_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/debounce"
)

func makeEntry(method string) capture.Entry {
	return capture.Entry{Method: method}
}

func TestAdd_CoalescesEntriesWithinWindow(t *testing.T) {
	var mu sync.Mutex
	var got [][]capture.Entry

	d := debounce.New(50*time.Millisecond, func(entries []capture.Entry) {
		mu.Lock()
		got = append(got, entries)
		mu.Unlock()
	})

	d.Add(makeEntry("/svc/A"))
	d.Add(makeEntry("/svc/B"))
	d.Add(makeEntry("/svc/C"))

	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(got) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(got))
	}
	if len(got[0]) != 3 {
		t.Fatalf("expected 3 entries in batch, got %d", len(got[0]))
	}
}

func TestAdd_ResetsTimerOnEachCall(t *testing.T) {
	fired := make(chan struct{}, 1)
	d := debounce.New(60*time.Millisecond, func(entries []capture.Entry) {
		fired <- struct{}{}
	})

	// Repeatedly add entries to keep resetting the timer.
	for i := 0; i < 5; i++ {
		d.Add(makeEntry("/svc/Loop"))
		time.Sleep(20 * time.Millisecond)
	}

	select {
	case <-fired:
		t.Fatal("handler fired too early")
	case <-time.After(10 * time.Millisecond):
	}

	// Now wait for the window to expire.
	time.Sleep(100 * time.Millisecond)
	select {
	case <-fired:
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("handler did not fire after window elapsed")
	}
}

func TestFlush_ImmediatelyInvokesHandler(t *testing.T) {
	called := make(chan []capture.Entry, 1)
	d := debounce.New(5*time.Second, func(entries []capture.Entry) {
		called <- entries
	})

	d.Add(makeEntry("/svc/Flush"))
	d.Flush()

	select {
	case batch := <-called:
		if len(batch) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(batch))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("handler not called after Flush")
	}
}

func TestFlush_NoPendingEntries_DoesNotCallHandler(t *testing.T) {
	called := false
	d := debounce.New(50*time.Millisecond, func(_ []capture.Entry) {
		called = true
	})
	d.Flush()
	if called {
		t.Fatal("handler should not be called when no entries are pending")
	}
}

func TestRun_FlushesOnContextCancel(t *testing.T) {
	called := make(chan struct{}, 1)
	d := debounce.New(5*time.Second, func(entries []capture.Entry) {
		if len(entries) > 0 {
			called <- struct{}{}
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	d.Run(ctx)
	d.Add(makeEntry("/svc/Ctx"))
	cancel()

	select {
	case <-called:
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("handler not called after context cancel")
	}
}
