// Package sampler implements probabilistic sampling for captured gRPC entries.
//
// A Sampler is configured with a rate in [0.0, 1.0] that controls what
// fraction of traffic is retained. A rate of 1.0 keeps every entry; a rate
// of 0.0 drops all entries.
//
// Example usage:
//
//	s := sampler.New(0.1, rand.NewSource(time.Now().UnixNano()))
//	kept := s.Filter(store.List())
//
// The sampling rate can be adjusted at runtime via SetRate, making it
// straightforward to dial traffic volume up or down without restarting the
// monitoring session.
package sampler
