package pipeline_test

import (
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/pipeline"
)

func entries(methods ...string) []capture.Entry {
	out := make([]capture.Entry, len(methods))
	for i, m := range methods {
		out[i] = capture.Entry{Method: m}
	}
	return out
}

func TestRun_NoSteps_ReturnsInput(t *testing.T) {
	p := pipeline.New()
	in := entries("/svc/A", "/svc/B")
	got := p.Run(in)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestRun_SingleStep_FiltersEntries(t *testing.T) {
	keepA := func(es []capture.Entry) []capture.Entry {
		var out []capture.Entry
		for _, e := range es {
			if e.Method == "/svc/A" {
				out = append(out, e)
			}
		}
		return out
	}
	p := pipeline.New(keepA)
	got := p.Run(entries("/svc/A", "/svc/B", "/svc/A"))
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestRun_MultipleSteps_AppliedInOrder(t *testing.T) {
	var order []int
	step := func(n int) pipeline.Processor {
		return func(es []capture.Entry) []capture.Entry {
			order = append(order, n)
			return es
		}
	}
	p := pipeline.New(step(1), step(2), step(3))
	p.Run(entries("/svc/X"))
	for i, v := range order {
		if v != i+1 {
			t.Fatalf("expected step %d at position %d, got %d", i+1, i, v)
		}
	}
}

func TestAdd_AppendsStep(t *testing.T) {
	called := false
	p := pipeline.New()
	p.Add(func(es []capture.Entry) []capture.Entry {
		called = true
		return es
	})
	p.Run(entries("/svc/Z"))
	if !called {
		t.Fatal("expected added step to be called")
	}
}
