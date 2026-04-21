package chain_test

import (
	"errors"
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/chain"
)

func sampleEntry() capture.Entry {
	return capture.Entry{Method: "/pkg.Service/Method"}
}

func TestNew_EmptyChain_CallsTerminal(t *testing.T) {
	c := chain.New()
	called := false
	h := c.Then(func(e capture.Entry) error {
		called = true
		return nil
	})
	if err := h(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected terminal handler to be called")
	}
}

func TestUse_MiddlewareExecutedInOrder(t *testing.T) {
	var order []int
	mk := func(n int) chain.Middleware {
		return func(next chain.Handler) chain.Handler {
			return func(e capture.Entry) error {
				order = append(order, n)
				return next(e)
			}
		}
	}

	c := chain.New(mk(1), mk(2))
	c.Use(mk(3))

	_ = c.Then(func(capture.Entry) error { return nil })(sampleEntry())

	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestThen_NilTerminal_DoesNotPanic(t *testing.T) {
	c := chain.New()
	h := c.Then(nil)
	if err := h(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_PropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	c := chain.New()
	err := c.Run(sampleEntry(), func(capture.Entry) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestMiddleware_CanShortCircuit(t *testing.T) {
	sentinel := errors.New("short-circuit")
	guard := chain.Middleware(func(next chain.Handler) chain.Handler {
		return func(e capture.Entry) error {
			return sentinel
		}
	})

	c := chain.New(guard)
	terminalCalled := false
	err := c.Run(sampleEntry(), func(capture.Entry) error {
		terminalCalled = true
		return nil
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
	if terminalCalled {
		t.Fatal("terminal should not have been called")
	}
}
