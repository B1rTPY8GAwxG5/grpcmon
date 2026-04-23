// Package transform provides a composable pipeline for mutating captured
// gRPC entries before they are replayed, exported, or displayed.
package transform

import "github.com/example/grpcmon/internal/capture"

// Func is a function that receives an entry and returns a (possibly modified)
// copy. Returning a zero-value Entry signals that the entry should be dropped.
type Func func(capture.Entry) (capture.Entry, bool)

// Chain holds an ordered sequence of Funcs and applies them in order.
type Chain struct {
	steps []Func
}

// New returns an empty Chain.
func New() *Chain {
	return &Chain{}
}

// Add appends a Func to the chain.
func (c *Chain) Add(f Func) *Chain {
	c.steps = append(c.steps, f)
	return c
}

// Apply runs all Funcs against entry. If any step signals a drop the entry is
// excluded from the result and Apply returns false.
func (c *Chain) Apply(e capture.Entry) (capture.Entry, bool) {
	for _, step := range c.steps {
		var ok bool
		e, ok = step(e)
		if !ok {
			return capture.Entry{}, false
		}
	}
	return e, true
}

// ApplyAll transforms a slice of entries, omitting any that are dropped by the
// chain.
func (c *Chain) ApplyAll(entries []capture.Entry) []capture.Entry {
	out := make([]capture.Entry, 0, len(entries))
	for _, e := range entries {
		if transformed, ok := c.Apply(e); ok {
			out = append(out, transformed)
		}
	}
	return out
}

// SetMethod returns a Func that replaces the Method field of every entry.
func SetMethod(method string) Func {
	return func(e capture.Entry) (capture.Entry, bool) {
		e.Method = method
		return e, true
	}
}

// DropErrors returns a Func that drops entries whose StatusCode is non-zero.
func DropErrors() Func {
	return func(e capture.Entry) (capture.Entry, bool) {
		if e.StatusCode != 0 {
			return capture.Entry{}, false
		}
		return e, true
	}
}

// OverrideTarget returns a Func that replaces the Target field of every entry.
func OverrideTarget(target string) Func {
	return func(e capture.Entry) (capture.Entry, bool) {
		e.Target = target
		return e, true
	}
}
