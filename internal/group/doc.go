// Package group provides utilities for grouping captured gRPC entries
// by an arbitrary key function such as method name or status code.
//
// Usage:
//
//	s := group.New(group.ByMethod)
//	groups := s.Apply(store.List())
//	for _, g := range groups {
//		fmt.Println(g.Key, len(g.Entries))
//	}
package group
