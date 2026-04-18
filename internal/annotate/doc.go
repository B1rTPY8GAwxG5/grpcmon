// Package annotate allows developers to attach free-form text notes to any
// captured gRPC entry by its unique ID.
//
// Notes are stored in memory and can be persisted to disk as JSON via the
// Save and LoadInto helpers, making them easy to share alongside snapshot
// files.
//
// Example usage:
//
//	s := annotate.New()
//	_ = s.Set(entryID, "retried after timeout – see issue #42")
//	note, _ := s.Get(entryID)
//	fmt.Println(note)
package annotate
