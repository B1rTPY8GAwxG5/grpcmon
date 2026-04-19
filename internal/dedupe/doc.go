// Package dedupe provides deduplication of captured gRPC entries.
//
// Entries are fingerprinted using a SHA-256 hash of their method and
// request payload. Identical method+request pairs are considered
// duplicates and can be filtered out before storage or replay.
//
// Example usage:
//
//	ds := dedupe.New()
//	unique := ds.Filter(store.List())
package dedupe
