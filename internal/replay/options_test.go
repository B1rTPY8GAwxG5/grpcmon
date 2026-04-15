package replay

import (
	"testing"
	"time"
)

func TestDefaultOptions_Values(t *testing.T) {
	opts := DefaultOptions()

	if opts.DelayBetween != 0 {
		t.Errorf("DelayBetween: got %v, want 0", opts.DelayBetween)
	}
	if opts.TimeoutPerRequest != 10*time.Second {
		t.Errorf("TimeoutPerRequest: got %v, want 10s", opts.TimeoutPerRequest)
	}
}

func TestOptions_CustomValues(t *testing.T) {
	opts := Options{
		DelayBetween:      100 * time.Millisecond,
		TimeoutPerRequest: 5 * time.Second,
	}

	if opts.DelayBetween != 100*time.Millisecond {
		t.Errorf("DelayBetween: got %v, want 100ms", opts.DelayBetween)
	}
	if opts.TimeoutPerRequest != 5*time.Second {
		t.Errorf("TimeoutPerRequest: got %v, want 5s", opts.TimeoutPerRequest)
	}
}
