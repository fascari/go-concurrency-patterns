//go:build mutexcopy

package mutexcopy

import "sync"

type (
	Store struct {
		mu       sync.Mutex
		balances map[string]int64
	}
)

func New() *Store {
	return new(Store{
		balances: make(map[string]int64),
	})
}

// DemonstrateCopiedLock exists so go vet reports the copied mutex.
func DemonstrateCopiedLock() {
	original := New()
	copied := *original
	copied.mu.Lock()
}
