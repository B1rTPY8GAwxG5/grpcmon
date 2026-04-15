package capture

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
)

var counter uint64

// NewID generates a unique entry ID combining a random prefix and a
// monotonically increasing counter for ordering guarantees.
func NewID() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		// Fallback to counter-only ID if crypto/rand is unavailable.
		return fmt.Sprintf("%016x", atomic.AddUint64(&counter, 1))
	}
	seq := atomic.AddUint64(&counter, 1)
	return fmt.Sprintf("%s-%08x", hex.EncodeToString(b), seq)
}
