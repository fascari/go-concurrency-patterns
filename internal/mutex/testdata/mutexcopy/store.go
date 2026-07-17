package mutexcopy

import "sync"

type Store struct {
	mu       sync.Mutex
	balances map[string]int64
}

func New() *Store {
	return new(Store{
		balances: make(map[string]int64),
	})
}

// DemonstrateCopiedLock exists solely so that go vet detects the copied lock
//
//goland:noinspection ALL
func DemonstrateCopiedLock() {
	original := New()
	copy := *original
	copy.mu.Lock()
}
