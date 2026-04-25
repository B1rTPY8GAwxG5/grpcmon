// Package budget provides an error-budget tracker for gRPC methods.
// It tracks the ratio of successful calls against a configurable SLO target
// and reports whether the budget has been exhausted.
package budget

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
)

// Budget tracks the error budget for a single method.
type Budget struct {
	mu       sync.Mutex
	target   float64 // e.g. 0.99 for 99% success rate
	total    int
	success  int
}

// Store holds budgets keyed by method name.
type Store struct {
	mu      sync.Mutex
	budgets map[string]*Budget
	target  float64
}

// New returns a Store where every method shares the given SLO target (0–1).
func New(target float64) *Store {
	if target <= 0 || target > 1 {
		target = 0.99
	}
	return &Store{
		budgets: make(map[string]*Budget),
		target:  target,
	}
}

// Record registers a call result for the given method.
func (s *Store) Record(method string, code codes.Code) {
	s.mu.Lock()
	b, ok := s.budgets[method]
	if !ok {
		b = &Budget{target: s.target}
		s.budgets[method] = b
	}
	s.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()
	b.total++
	if code == codes.OK {
		b.success++
	}
}

// Remaining returns the fraction of error budget remaining for the method.
// A value >= 0 means budget is intact; negative means it is exhausted.
func (s *Store) Remaining(method string) float64 {
	s.mu.Lock()
	b, ok := s.budgets[method]
	s.mu.Unlock()
	if !ok || b.total == 0 {
		return 1.0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	actualRate := float64(b.success) / float64(b.total)
	// remaining = (actual - target) / (1 - target)
	if b.target >= 1 {
		return 0
	}
	return (actualRate - b.target) / (1 - b.target)
}

// Exhausted reports whether the error budget for method has been spent.
func (s *Store) Exhausted(method string) bool {
	return s.Remaining(method) < 0
}

// Summary returns a human-readable budget summary for method.
func (s *Store) Summary(method string) string {
	r := s.Remaining(method)
	status := "OK"
	if r < 0 {
		status = "EXHAUSTED"
	}
	return fmt.Sprintf("method=%s remaining=%.2f%% status=%s", method, r*100, status)
}
