// Package routing provides method-based routing for captured gRPC entries,
// allowing handlers to be registered per gRPC method path.
package routing

import (
	"fmt"
	"sync"

	"github.com/grpcmon/internal/capture"
)

// Handler is a function that processes a captured entry.
type Handler func(entry capture.Entry)

// Router dispatches captured entries to registered handlers based on gRPC method.
type Router struct {
	mu       sync.RWMutex
	routes   map[string]Handler
	fallback Handler
}

// New creates a new Router with an optional fallback handler
// that is called when no specific route matches.
func New(fallback Handler) *Router {
	return &Router{
		routes:   make(map[string]Handler),
		fallback: fallback,
	}
}

// Register associates a Handler with the given gRPC method path.
// Registering the same method twice overwrites the previous handler.
func (r *Router) Register(method string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[method] = h
}

// Deregister removes the handler for the given method, if any.
func (r *Router) Deregister(method string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.routes, method)
}

// Dispatch routes the entry to the matching handler.
// If no handler is registered for the method, the fallback is called.
// If no fallback is set, Dispatch returns an error.
func (r *Router) Dispatch(entry capture.Entry) error {
	r.mu.RLock()
	h, ok := r.routes[entry.Method]
	fb := r.fallback
	r.mu.RUnlock()

	if ok {
		h(entry)
		return nil
	}
	if fb != nil {
		fb(entry)
		return nil
	}
	return fmt.Errorf("routing: no handler registered for method %q", entry.Method)
}

// Methods returns a sorted snapshot of all registered method paths.
func (r *Router) Methods() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.routes))
	for m := range r.routes {
		out = append(out, m)
	}
	return out
}
