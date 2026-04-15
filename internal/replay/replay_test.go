package replay

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/grpcmon/internal/capture"
)

// startNoopServer spins up a minimal gRPC server that rejects all calls.
func startNoopServer(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	go srv.Serve(lis) //nolint:errcheck
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

func TestNew_ConnectsSuccessfully(t *testing.T) {
	addr := startNoopServer(t)
	r, err := New(addr, DefaultOptions())
	if err != nil {
		t.Fatalf("New: unexpected error: %v", err)
	}
	defer r.Close()
}

func TestReplayEntry_ReturnsResult(t *testing.T) {
	addr := startNoopServer(t)
	r, err := New(addr, DefaultOptions())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	e := capture.Entry{
		ID:     "test-id",
		Method: "/pkg.Service/Method",
	}

	res := r.ReplayEntry(context.Background(), e)
	if res.EntryID != e.ID {
		t.Errorf("EntryID: got %q, want %q", res.EntryID, e.ID)
	}
	if res.Method != e.Method {
		t.Errorf("Method: got %q, want %q", res.Method, e.Method)
	}
	// The server has no registered service so an error is expected.
	if res.Err == nil {
		t.Error("expected error from unregistered method, got nil")
	}
}

func TestReplayAll_ClosesChannel(t *testing.T) {
	addr := startNoopServer(t)
	r, err := New(addr, DefaultOptions())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	entries := []capture.Entry{
		{ID: "a", Method: "/svc/A"},
		{ID: "b", Method: "/svc/B"},
	}

	ch := r.ReplayAll(context.Background(), entries)
	var results []Result
	for res := range ch {
		results = append(results, res)
	}
	if len(results) != len(entries) {
		t.Errorf("got %d results, want %d", len(results), len(entries))
	}
}

func TestReplayAll_RespectsContextCancel(t *testing.T) {
	addr := startNoopServer(t)
	opts := DefaultOptions()
	opts.DelayBetween = 50 * time.Millisecond
	r, err := New(addr, opts)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	entries := make([]capture.Entry, 10)
	for i := range entries {
		entries[i] = capture.Entry{ID: "x", Method: "/svc/X"}
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := r.ReplayAll(ctx, entries)
	// Cancel after first result.
	<-ch
	cancel()
	// Drain; must not block forever.
	done := make(chan struct{})
	go func() { for range ch {}; close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("ReplayAll did not respect context cancellation")
	}
}
