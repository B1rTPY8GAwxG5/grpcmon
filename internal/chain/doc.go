// Package chain provides a composable middleware chain for processing
// captured gRPC entries.
//
// Middleware functions wrap a Handler, enabling cross-cutting concerns such as
// logging, filtering, rate-limiting, or masking to be applied in a declarative
// pipeline without modifying core business logic.
//
// Example usage:
//
//	c := chain.New(
//		loggingMiddleware,
//		maskingMiddleware,
//	)
//	c.Use(rateLimitMiddleware)
//
//	h := c.Then(func(e capture.Entry) error {
//		// handle entry
//		return nil
//	})
//
//	if err := h(entry); err != nil {
//		// handle error
//	}
package chain
