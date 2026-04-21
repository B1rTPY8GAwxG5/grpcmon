// Package chain provides a composable middleware chain for processing
// captured gRPC entries through a series of transformation steps.
package chain

import "github.com/grpcmon/internal/capture"

// Handler is a function that processes a single capture entry.
type Handler func(entry capture.Entry) error

// Middleware wraps a Handler, allowing pre/post processing.
type Middleware func(next Handler) Handler

// Chain holds an ordered list of middleware and a terminal handler.
type Chain struct {
	middlewares []Middleware
}

// New returns an empty Chain.
func New(mws ...Middleware) *Chain {
	c := &Chain{}
	for _, mw := range mws {
		c.middlewares = append(c.middlewares, mw)
	}
	return c
}

// Use appends one or more middleware to the chain.
func (c *Chain) Use(mws ...Middleware) {
	c.middlewares = append(c.middlewares, mws...)
}

// Then builds the final Handler by wrapping h with all middleware in order.
// The first middleware added is the outermost wrapper.
func (c *Chain) Then(h Handler) Handler {
	if h == nil {
		h = func(capture.Entry) error { return nil }
	}
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}

// Run is a convenience method that builds the chain and immediately invokes
// it with the provided entry.
func (c *Chain) Run(entry capture.Entry, h Handler) error {
	return c.Then(h)(entry)
}
