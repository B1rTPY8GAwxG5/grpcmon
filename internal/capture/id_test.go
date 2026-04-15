package capture

import (
	"strings"
	"sync"
	"testing"
)

func TestNewID_Format(t *testing.T) {
	id := NewID()
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts separated by '-', got %q", id)
	}
	if len(parts[0]) != 8 {
		t.Errorf("expected 8-char hex prefix, got %q", parts[0])
	}
	if len(parts[1]) != 8 {
		t.Errorf("expected 8-char hex counter, got %q", parts[1])
	}
}

func TestNewID_Uniqueness(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		id := NewID()
		if _, exists := seen[id]; exists {
			t.Fatalf("duplicate ID generated: %s", id)
		}
		seen[id] = struct{}{}
	}
}

func TestNewID_ConcurrentSafety(t *testing.T) {
	const goroutines = 50
	const perGoroutine = 100

	var mu sync.Mutex
	seen := make(map[string]struct{}, goroutines*perGoroutine)
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			ids := make([]string, perGoroutine)
			for j := 0; j < perGoroutine; j++ {
				ids[j] = NewID()
			}
			mu.Lock()
			for _, id := range ids {
				if _, exists := seen[id]; exists {
					t.Errorf("duplicate ID in concurrent test: %s", id)
				}
				seen[id] = struct{}{}
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
}
