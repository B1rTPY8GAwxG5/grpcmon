package transform_test

import (
	"testing"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/transform"
)

func baseEntry() capture.Entry {
	return capture.Entry{
		ID:         "id-1",
		Method:     "/svc/Method",
		Target:     "localhost:50051",
		StatusCode: 0,
	}
}

func TestApply_PassesThroughWhenNoSteps(t *testing.T) {
	c := transform.New()
	e := baseEntry()
	out, ok := c.Apply(e)
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if out.Method != e.Method {
		t.Errorf("method changed unexpectedly: got %q", out.Method)
	}
}

func TestSetMethod_ChangesMethod(t *testing.T) {
	c := transform.New().Add(transform.SetMethod("/new/Method"))
	out, ok := c.Apply(baseEntry())
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if out.Method != "/new/Method" {
		t.Errorf("expected /new/Method, got %q", out.Method)
	}
}

func TestDropErrors_DropsNonZeroStatus(t *testing.T) {
	c := transform.New().Add(transform.DropErrors())
	e := baseEntry()
	e.StatusCode = 2
	_, ok := c.Apply(e)
	if ok {
		t.Fatal("expected entry to be dropped")
	}
}

func TestDropErrors_KeepsOKEntry(t *testing.T) {
	c := transform.New().Add(transform.DropErrors())
	_, ok := c.Apply(baseEntry())
	if !ok {
		t.Fatal("expected entry to be kept")
	}
}

func TestOverrideTarget_ReplacesTarget(t *testing.T) {
	c := transform.New().Add(transform.OverrideTarget("newhost:9090"))
	out, ok := c.Apply(baseEntry())
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if out.Target != "newhost:9090" {
		t.Errorf("expected newhost:9090, got %q", out.Target)
	}
}

func TestApplyAll_FiltersAndTransforms(t *testing.T) {
	c := transform.New().
		Add(transform.DropErrors()).
		Add(transform.OverrideTarget("prod:443"))

	entries := []capture.Entry{
		{ID: "1", StatusCode: 0, Target: "old"},
		{ID: "2", StatusCode: 3, Target: "old"},
		{ID: "3", StatusCode: 0, Target: "old"},
	}

	out := c.ApplyAll(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Target != "prod:443" {
			t.Errorf("expected target prod:443, got %q", e.Target)
		}
	}
}

func TestChain_StepsAppliedInOrder(t *testing.T) {
	var order []string
	step := func(label string) transform.Func {
		return func(e capture.Entry) (capture.Entry, bool) {
			order = append(order, label)
			return e, true
		}
	}

	c := transform.New().Add(step("a")).Add(step("b")).Add(step("c"))
	c.Apply(baseEntry()) //nolint:errcheck

	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Errorf("unexpected step order: %v", order)
	}
}
