package routing_test

import (
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/routing"
)

func entry(method string) capture.Entry {
	return capture.Entry{Method: method}
}

func TestDispatch_CallsRegisteredHandler(t *testing.T) {
	r := routing.New(nil)
	var called string
	r.Register("/pkg.Service/Hello", func(e capture.Entry) {
		called = e.Method
	})
	if err := r.Dispatch(entry("/pkg.Service/Hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != "/pkg.Service/Hello" {
		t.Errorf("expected handler called with /pkg.Service/Hello, got %q", called)
	}
}

func TestDispatch_FallsBackWhenNoMatch(t *testing.T) {
	var fallbackCalled bool
	r := routing.New(func(e capture.Entry) { fallbackCalled = true })
	if err := r.Dispatch(entry("/pkg.Service/Unknown")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fallbackCalled {
		t.Error("expected fallback to be called")
	}
}

func TestDispatch_ErrorWhenNoHandlerAndNoFallback(t *testing.T) {
	r := routing.New(nil)
	err := r.Dispatch(entry("/pkg.Service/Missing"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeregister_RemovesHandler(t *testing.T) {
	r := routing.New(nil)
	r.Register("/pkg.Service/Hello", func(e capture.Entry) {})
	r.Deregister("/pkg.Service/Hello")
	err := r.Dispatch(entry("/pkg.Service/Hello"))
	if err == nil {
		t.Fatal("expected error after deregister, got nil")
	}
}

func TestMethods_ReturnsRegistered(t *testing.T) {
	r := routing.New(nil)
	r.Register("/a", func(e capture.Entry) {})
	r.Register("/b", func(e capture.Entry) {})
	ms := r.Methods()
	if len(ms) != 2 {
		t.Errorf("expected 2 methods, got %d", len(ms))
	}
}

func TestRegister_OverwritesPreviousHandler(t *testing.T) {
	r := routing.New(nil)
	calls := 0
	r.Register("/svc/M", func(e capture.Entry) { calls++ })
	r.Register("/svc/M", func(e capture.Entry) { calls += 10 })
	_ = r.Dispatch(entry("/svc/M"))
	if calls != 10 {
		t.Errorf("expected overwritten handler to be called (calls=10), got %d", calls)
	}
}
