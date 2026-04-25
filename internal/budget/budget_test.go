package budget_test

import (
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/example/grpcmon/internal/budget"
)

func TestNew_DefaultTarget(t *testing.T) {
	s := budget.New(0) // invalid → defaults to 0.99
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestRecord_And_Remaining_AllSuccess(t *testing.T) {
	s := budget.New(0.99)
	for i := 0; i < 100; i++ {
		s.Record("/svc/Method", codes.OK)
	}
	if s.Exhausted("/svc/Method") {
		t.Error("budget should not be exhausted when all calls succeed")
	}
	if r := s.Remaining("/svc/Method"); r <= 0 {
		t.Errorf("expected positive remaining, got %.4f", r)
	}
}

func TestRecord_And_Remaining_AllErrors(t *testing.T) {
	s := budget.New(0.99)
	for i := 0; i < 100; i++ {
		s.Record("/svc/Method", codes.Internal)
	}
	if !s.Exhausted("/svc/Method") {
		t.Error("budget should be exhausted when all calls fail")
	}
}

func TestRemaining_UnknownMethod(t *testing.T) {
	s := budget.New(0.99)
	if r := s.Remaining("/unknown/Method"); r != 1.0 {
		t.Errorf("expected 1.0 for unknown method, got %.4f", r)
	}
}

func TestExhausted_AtBoundary(t *testing.T) {
	s := budget.New(0.90)
	// 90 OK + 10 errors → success rate exactly 0.90 = target → remaining == 0
	for i := 0; i < 90; i++ {
		s.Record("/svc/M", codes.OK)
	}
	for i := 0; i < 10; i++ {
		s.Record("/svc/M", codes.Unavailable)
	}
	// remaining should be 0, not negative
	if s.Exhausted("/svc/M") {
		t.Error("budget should not be exhausted exactly at target")
	}
}

func TestSummary_ContainsMethod(t *testing.T) {
	s := budget.New(0.99)
	s.Record("/svc/Hello", codes.OK)
	sum := s.Summary("/svc/Hello")
	if sum == "" {
		t.Fatal("expected non-empty summary")
	}
	for _, want := range []string{"method=/svc/Hello", "remaining=", "status="} {
		if !contains(sum, want) {
			t.Errorf("summary missing %q: %s", want, sum)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := range s {
		if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
