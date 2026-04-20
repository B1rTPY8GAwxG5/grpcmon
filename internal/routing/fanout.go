package routing

import "github.com/grpcmon/internal/capture"

// Fanout distributes a single entry to multiple handlers in order.
// All handlers are called regardless of individual outcomes.
type Fanout struct {
	handlers []Handler
}

// NewFanout creates a Fanout that will call each of the provided handlers
// when its Handle method is invoked.
func NewFanout(handlers ...Handler) *Fanout {
	copied := make([]Handler, len(handlers))
	copy(copied, handlers)
	return &Fanout{handlers: copied}
}

// Add appends an additional handler to the fanout.
func (f *Fanout) Add(h Handler) {
	f.handlers = append(f.handlers, h)
}

// Handle calls every registered handler with the given entry.
func (f *Fanout) Handle(entry capture.Entry) {
	for _, h := range f.handlers {
		h(entry)
	}
}

// AsHandler returns the Fanout as a routing.Handler for use with Router.Register.
func (f *Fanout) AsHandler() Handler {
	return f.Handle
}
