// Package tag provides lightweight tagging support for captured gRPC entries.
//
// Tags are arbitrary string labels that can be associated with one or more
// entry IDs, enabling fast lookup and filtering without modifying the
// underlying capture.Entry struct.
//
// Example usage:
//
//	ts := tag.New()
//	ts.Add(entryID, "slow", "retried")
//	matches := ts.Filter("slow", allEntries)
package tag
