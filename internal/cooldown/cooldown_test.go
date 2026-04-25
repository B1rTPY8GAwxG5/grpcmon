package cooldown

import (
	"testing"
	"time"
)

func TestNew_DefaultsToOneSecond(t *testing.T) {
	c := New(0)
	if c.interval != time.Second {
		t.Fatalf("expected 1s interval, got %v", c.interval)
	}
}

func TestAllow_FirstCallSucceeds(t *testing.T) {
	c := New(100 * time.Millisecond)
	if !c.Allow("/svc/Method") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_SecondCallWithinIntervalBlocked(t *testing.T) {
	c := New(500 * time.Millisecond)
	c.Allow("/svc/Method")
	if c.Allow("/svc/Method") {
		t.Fatal("expected second Allow within interval to return false")
	}
}

func TestAllow_AllowedAfterIntervalExpires(t *testing.T) {
	c := New(20 * time.Millisecond)
	c.Allow("/svc/Method")
	time.Sleep(30 * time.Millisecond)
	if !c.Allow("/svc/Method") {
		t.Fatal("expected Allow to return true after interval elapsed")
	}
}

func TestAllow_IndependentPerMethod(t *testing.T) {
	c := New(500 * time.Millisecond)
	c.Allow("/svc/A")
	if !c.Allow("/svc/B") {
		t.Fatal("expected different method to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	c := New(500 * time.Millisecond)
	c.Allow("/svc/Method")
	c.Reset("/svc/Method")
	if !c.Allow("/svc/Method") {
		t.Fatal("expected Allow to succeed after Reset")
	}
}

func TestResetAll_ClearsAllMethods(t *testing.T) {
	c := New(500 * time.Millisecond)
	c.Allow("/svc/A")
	c.Allow("/svc/B")
	c.ResetAll()
	if !c.Allow("/svc/A") || !c.Allow("/svc/B") {
		t.Fatal("expected all methods to be allowed after ResetAll")
	}
}

func TestRemaining_ZeroForUnknownMethod(t *testing.T) {
	c := New(500 * time.Millisecond)
	if r := c.Remaining("/svc/Unknown"); r != 0 {
		t.Fatalf("expected 0 remaining for unknown method, got %v", r)
	}
}

func TestRemaining_PositiveWithinInterval(t *testing.T) {
	c := New(500 * time.Millisecond)
	c.Allow("/svc/Method")
	r := c.Remaining("/svc/Method")
	if r <= 0 || r > 500*time.Millisecond {
		t.Fatalf("expected remaining in (0, 500ms], got %v", r)
	}
}

func TestRemaining_ZeroAfterIntervalExpires(t *testing.T) {
	c := New(20 * time.Millisecond)
	c.Allow("/svc/Method")
	time.Sleep(30 * time.Millisecond)
	if r := c.Remaining("/svc/Method"); r != 0 {
		t.Fatalf("expected 0 remaining after interval, got %v", r)
	}
}
