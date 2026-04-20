package routing_test

import (
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/routing"
)

func TestFanout_CallsAllHandlers(t *testing.T) {
	calls := make([]int, 3)
	f := routing.NewFanout(
		func(e capture.Entry) { calls[0]++ },
		func(e capture.Entry) { calls[1]++ },
		func(e capture.Entry) { calls[2]++ },
	)
	f.Handle(capture.Entry{Method: "/svc/M"})
	for i, c := range calls {
		if c != 1 {
			t.Errorf("handler %d: expected 1 call, got %d", i, c)
		}
	}
}

func TestFanout_Add_AppendsDynamically(t *testing.T) {
	calls := 0
	f := routing.NewFanout()
	f.Add(func(e capture.Entry) { calls++ })
	f.Add(func(e capture.Entry) { calls++ })
	f.Handle(capture.Entry{})
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestFanout_AsHandler_IntegratesWithRouter(t *testing.T) {
	received := ""
	f := routing.NewFanout(func(e capture.Entry) { received = e.Method })

	r := routing.New(nil)
	r.Register("/svc/Ping", f.AsHandler())

	if err := r.Dispatch(capture.Entry{Method: "/svc/Ping"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received != "/svc/Ping" {
		t.Errorf("expected /svc/Ping, got %q", received)
	}
}

func TestFanout_EmptyFanout_DoesNotPanic(t *testing.T) {
	f := routing.NewFanout()
	f.Handle(capture.Entry{Method: "/svc/M"})
}
